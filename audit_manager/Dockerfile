FROM python:3.11-slim

WORKDIR /app

# Instalar dependências do sistema
RUN apt-get update && apt-get install -y \
    gcc \
    && rm -rf /var/lib/apt/lists/*

# Copiar arquivos de dependências
COPY requirements.txt .

# Instalar dependências Python
RUN pip install --no-cache-dir -r requirements.txt

# Copiar código da aplicação
COPY main.py .

# Criar usuário não-root
RUN useradd -m -u 1000 auditmanager && chown -R auditmanager:auditmanager /app
USER auditmanager

# Comando padrão
CMD ["python", "main.py"]

