import os
import json
import logging
import signal
import sys
from datetime import datetime
from typing import Dict, Any

import pika
from opensearchpy import OpenSearch, RequestsHttpConnection
from opensearchpy.exceptions import OpenSearchException
from dotenv import load_dotenv

load_dotenv()

# Configurar logging
logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s'
)
logger = logging.getLogger(__name__)


class AuditManager:
    """Gerenciador de auditoria que consome eventos do RabbitMQ e armazena no OpenSearch"""
    
    def __init__(self):
        """Inicializa conexões com RabbitMQ e OpenSearch"""
        self.rabbitmq_connection = None
        self.rabbitmq_channel = None
        self.opensearch_client = None
        self.running = True
        
        # Carregar variáveis de ambiente
        # Configurações RabbitMQ
        self.host = os.getenv('RABBITMQ_HOST')
        self.port = int(os.getenv('RABBITMQ_PORT'))
        self.username = os.getenv('RABBITMQ_USERNAME')
        self.password = os.getenv('RABBITMQ_PASSWORD')
        self.queue = os.getenv('RABBITMQ_AUDIT_QUEUE')
        self.rabbitmq_exchange = os.getenv('RABBITMQ_EXCHANGE')
        self.rabbitmq_routing_key = os.getenv('RABBITMQ_ROUTING_KEY')
        
        # Configurações OpenSearch
        self.opensearch_host = os.getenv('OPENSEARCH_HOST')
        self.opensearch_port = int(os.getenv('OPENSEARCH_PORT'))
        self.opensearch_username = os.getenv('OPENSEARCH_USERNAME')
        self.opensearch_password = os.getenv('OPENSEARCH_PASSWORD')
        self.opensearch_use_ssl = os.getenv('OPENSEARCH_USE_SSL').lower() == 'true'
        self.opensearch_verify_certs = os.getenv('OPENSEARCH_VERIFY_CERTS').lower() == 'true'
        self.opensearch_index = os.getenv('OPENSEARCH_INDEX')
        
        # Registrar handlers de sinal para shutdown graceful
        signal.signal(signal.SIGINT, self._signal_handler)
        signal.signal(signal.SIGTERM, self._signal_handler)
    
    def _signal_handler(self, signum, frame):
        logger.info(f"Recebido sinal {signum}, encerrando...")
        self.running = False
    
    def connect_rabbitmq(self):
        """Conecta ao RabbitMQ"""
        try:
            credentials = pika.PlainCredentials(
                self.username,
                self.password
            )
            
            parameters = pika.ConnectionParameters(
                host=self.host,
                port=self.port,
                credentials=credentials,
                heartbeat=600,
                blocked_connection_timeout=300
            )
            
            self.rabbitmq_connection = pika.BlockingConnection(parameters)
            self.rabbitmq_channel = self.rabbitmq_connection.channel()
            
            self.rabbitmq_channel.queue_declare(
                queue=self.queue,
                durable=False
            )
            
            logger.info(f"Conectado ao RabbitMQ em {self.host}:{self.port}")
            logger.info(f"Consumindo fila: {self.queue}")
            
        except Exception as e:
            logger.error(f"Erro ao conectar ao RabbitMQ: {e}")
            raise
    
    def connect_opensearch(self):
        try:
            self.opensearch_client = OpenSearch(
                hosts=[{
                    'host': self.opensearch_host,
                    'port': self.opensearch_port
                }],
                http_auth=(self.opensearch_username, self.opensearch_password),
                use_ssl=self.opensearch_use_ssl,
                verify_certs=self.opensearch_verify_certs,
                ssl_show_warn=False,
                connection_class=RequestsHttpConnection
            )
            
            # Verificar conexão
            info = self.opensearch_client.info()
            logger.info(f"Conectado ao OpenSearch: {info['version']['number']}")
            
            # Criar índice se não existir
            self._ensure_index_exists()
            
        except Exception as e:
            logger.error(f"Erro ao conectar ao OpenSearch: {e}")
            raise
    
    def _ensure_index_exists(self):
        """Garante que o índice existe no OpenSearch"""
        try:
            if not self.opensearch_client.indices.exists(index=self.opensearch_index):
                # Criar índice com mapeamento
                index_body = {
                    "settings": {
                        "number_of_shards": 1,
                        "number_of_replicas": 1
                    },
                    "mappings": {
                        "properties": {
                            "vehicle_type": {"type": "keyword"},
                            "vehicle_uuid": {"type": "keyword"},
                            "structure_type": {"type": "keyword"},
                            "structure_uuid": {"type": "keyword"},
                            "timestamp": {"type": "date"},
                            "event": {"type": "keyword"},  
                            "result": {"type": "keyword"},  
                            "slot_number": {"type": "integer"},
                            "tower_id": {"type": "keyword"}, 
                            "event_type": {"type": "keyword"}, 
                            "indexed_at": {"type": "date"}
                        }
                    }
                }
                
                try:
                    self.opensearch_client.indices.create(
                        index=self.opensearch_index,
                        body=index_body
                    )
                except TypeError:
                    # Fallback para versões que não usam 'body'
                    self.opensearch_client.indices.create(
                        index=self.opensearch_index,
                        **index_body
                    )
                logger.info(f"Índice '{self.opensearch_index}' criado com sucesso")
            else:
                logger.info(f"Índice '{self.opensearch_index}' já existe")
                
        except OpenSearchException as e:
            logger.error(f"Erro ao criar/verificar índice: {e}")
            raise
    
    def process_audit_event(self, body: bytes) -> Dict[str, Any]:
        """
        Processa um evento de auditoria recebido do RabbitMQ
        
        Suport para dois tipos de eventos:
        1. Events (H/S -> RabbitMQ): eventos de chegada/partida com campo 'event' e 'tower_id'
        2. Requests (T -> RabbitMQ): requisições de slots com campo 'result'
        
        Args:
            body: Corpo da mensagem em bytes
            
        Returns:
            Dicionário com os dados do evento processado
        """
        try:
            event_data = json.loads(body.decode('utf-8'))
            
            common_fields = [
                'vehicle_type', 'vehicle_uuid', 'structure_type',
                'structure_uuid', 'timestamp', 'slot_number'
            ]
            
            for field in common_fields:
                if field not in event_data:
                    raise ValueError(f"Campo obrigatório ausente: {field}")
            
            # Determinar tipo de evento
            has_event = 'event' in event_data
            has_result = 'result' in event_data
            
            if has_event and has_result:
                raise ValueError("Evento não pode ter ambos 'event' e 'result'")
            elif has_event:
                # Evento do tipo Events (H/S -> RabbitMQ)
                event_data['event_type'] = 'event'
                if 'tower_id' not in event_data:
                    raise ValueError("Eventos do tipo 'event' devem ter 'tower_id'")
            elif has_result:
                # Evento do tipo Requests (T -> RabbitMQ)
                event_data['event_type'] = 'request'
            else:
                raise ValueError("Evento deve ter 'event' ou 'result'")
            
            if isinstance(event_data['timestamp'], (int, float)):
                event_data['timestamp'] = datetime.utcfromtimestamp(event_data['timestamp'])
            elif isinstance(event_data['timestamp'], str):
                try:
                    event_data['timestamp'] = datetime.fromisoformat(
                        event_data['timestamp'].replace('Z', '+00:00')
                    )
                except ValueError:
                    try:
                        event_data['timestamp'] = datetime.utcfromtimestamp(
                            float(event_data['timestamp'])
                        )
                    except (ValueError, TypeError):
                        raise ValueError(f"Formato de timestamp inválido: {event_data['timestamp']}")
            
            # Adicionar timestamp de indexação
            event_data['indexed_at'] = datetime.utcnow()
            
            logger.debug(f"Evento processado: {event_data}")
            return event_data
            
        except json.JSONDecodeError as e:
            logger.error(f"Erro ao decodificar JSON: {e}")
            raise
        except Exception as e:
            logger.error(f"Erro ao processar evento: {e}")
            raise
    
    def store_in_opensearch(self, event_data: Dict[str, Any]):
        """
        Armazena evento no OpenSearch
        
        Args:
            event_data: Dados do evento a serem armazenados
        """
        try:
            doc = event_data.copy()
            if isinstance(doc['timestamp'], datetime):
                doc['timestamp'] = doc['timestamp'].isoformat()
            if isinstance(doc['indexed_at'], datetime):
                doc['indexed_at'] = doc['indexed_at'].isoformat()
            
            # Indexar documento
            try:
                response = self.opensearch_client.index(
                    index=self.opensearch_index,
                    body=doc
                )
            except TypeError:
                # Fallback para versões que não usam 'body'
                response = self.opensearch_client.index(
                    index=self.opensearch_index,
                    **doc
                )
            
            event_info = ""
            if event_data.get('event_type') == 'event':
                event_info = f"Event: {event_data.get('event')}, Tower: {event_data.get('tower_id')}"
            else:
                event_info = f"Result: {event_data.get('result')}"
            
            logger.info(
                f"Evento armazenado no OpenSearch - "
                f"ID: {response['_id']}, "
                f"Type: {event_data.get('event_type')}, "
                f"Vehicle: {event_data['vehicle_uuid']}, "
                f"{event_info}"
            )
            
        except OpenSearchException as e:
            logger.error(f"Erro ao armazenar no OpenSearch: {e}")
            raise
    
    def on_message(self, channel, method, properties, body):
        """
        Callback chamado quando uma mensagem é recebida do RabbitMQ
        
        Args:
            channel: Canal do RabbitMQ
            method: Método de entrega
            properties: Propriedades da mensagem
            body: Corpo da mensagem
        """
        try:
            logger.info(f"Evento recebido: {body.decode('utf-8')}")
            
            event_data = self.process_audit_event(body)
            
            self.store_in_opensearch(event_data)
            
            channel.basic_ack(delivery_tag=method.delivery_tag)
            
        except Exception as e:
            logger.error(f"Erro ao processar mensagem: {e}")
            channel.basic_nack(
                delivery_tag=method.delivery_tag,
                requeue=False
            )
    
    def start_consuming(self):
        """Inicia o consumo de mensagens do RabbitMQ"""
        try:
            self.rabbitmq_channel.basic_qos(prefetch_count=1)
            
            self.rabbitmq_channel.basic_consume(
                queue=self.queue,
                on_message_callback=self.on_message,
                auto_ack=False
            )
            
            logger.info("Aguardando eventos de auditoria...")
            logger.info("Pressione CTRL+C para encerrar")
            
            while self.running:
                try:
                    self.rabbitmq_connection.process_data_events(time_limit=1)
                except KeyboardInterrupt:
                    break
            
        except Exception as e:
            logger.error(f"Erro ao consumir mensagens: {e}")
            raise
    
    def shutdown(self):
        """Encerra conexões de forma segura"""
        logger.info("Encerrando conexões...")
        
        if self.rabbitmq_channel and not self.rabbitmq_channel.is_closed:
            self.rabbitmq_channel.stop_consuming()
            self.rabbitmq_channel.close()
        
        if self.rabbitmq_connection and not self.rabbitmq_connection.is_closed:
            self.rabbitmq_connection.close()
        
        logger.info("Conexões encerradas")
    
    def run(self):
        """Executa o Audit Manager"""
        try:
            logger.info("Iniciando Audit Manager...")
            
            # Conectar aos serviços
            self.connect_rabbitmq()
            self.connect_opensearch()
            
            # Iniciar consumo
            self.start_consuming()
            
        except Exception as e:
            logger.error(f"Erro fatal: {e}")
            sys.exit(1)
        finally:
            self.shutdown()


def main():
    """Função principal"""
    manager = AuditManager()
    manager.run()


if __name__ == "__main__":
    main()