import { useState } from 'react';
import { BrowserRouter as Router, Routes, Route, Link, useNavigate } from 'react-router-dom';
import { motion } from 'framer-motion';
import './index.css'
import towerIcon from './assets/broca-de-petroleo.png';
import helicopterIcon from './assets/helicoptero.png';
import shipIcon from './assets/barco.png';
import platformIcon from './assets/poco-de-petroleo.png';
import centralIcon from './assets/plataforma-de-petroleo.png';
import towerLeaderIcon from './assets/tower-icon.png';

// MOCK_DATA (sem alterações)
const MOCK_TOWERS = [
  { id: 1, name: 'Tower Alpha', latitude: -23.55, longitude: -46.63, is_leader: false },
  { id: 1, name: 'Tower Alpha', latitude: -23.55, longitude: -46.63, is_leader: false },
  { id: 1, name: 'Tower Alpha', latitude: -23.55, longitude: -46.63, is_leader: false },
  { id: 1, name: 'Tower Alpha', latitude: -23.55, longitude: -46.63, is_leader: false },
  { id: 1, name: 'Tower Alpha', latitude: -23.55, longitude: -46.63, is_leader: false },
  { id: 1, name: 'Tower Alpha', latitude: -23.55, longitude: -46.63, is_leader: false },
  { id: 2, name: 'Tower Bravo', latitude: -22.9, longitude: -43.2, is_leader: true },
];
const MOCK_VEHICLES = [
  { id: 1, name: 'Helicopter A', type: 'Helicopter', latitude: -23.5, longitude: -46.6 },
  { id: 1, name: 'Helicopter A', type: 'Helicopter', latitude: -23.5, longitude: -46.6 },
  { id: 1, name: 'Helicopter A', type: 'Helicopter', latitude: -23.5, longitude: -46.6 },
  { id: 1, name: 'Helicopter A', type: 'Helicopter', latitude: -23.5, longitude: -46.6 },
  { id: 2, name: 'Ship B', type: 'Ship', latitude: -22.9, longitude: -43.2 },
  { id: 2, name: 'Ship B', type: 'Ship', latitude: -22.9, longitude: -43.2 },
  { id: 2, name: 'Ship B', type: 'Ship', latitude: -22.9, longitude: -43.2 },
  { id: 2, name: 'Ship B', type: 'Ship', latitude: -22.9, longitude: -43.2 },
];
const MOCK_STRUCTURES = [
  { id: 1, name: 'Platform X', latitude: -24.55, longitude: -46.62, type: 'Platform' },
  { id: 1, name: 'Platform X', latitude: -24.55, longitude: -46.62, type: 'Platform' },
  { id: 1, name: 'Platform X', latitude: -24.55, longitude: -46.62, type: 'Platform' },
  { id: 1, name: 'Platform X', latitude: -24.55, longitude: -46.62, type: 'Platform' },
  { id: 1, name: 'Platform X', latitude: -24.55, longitude: -46.62, type: 'Platform' },
  { id: 2, name: 'Central Y', latitude: -22.95, longitude: -43.25, type: 'Central' },
  { id: 2, name: 'Central Y', latitude: -22.95, longitude: -43.25, type: 'Central' },
  { id: 2, name: 'Central Y', latitude: -22.95, longitude: -43.25, type: 'Central' },
  { id: 2, name: 'Central Y', latitude: -22.95, longitude: -43.25, type: 'Central' },
  { id: 2, name: 'Central Y', latitude: -22.95, longitude: -43.25, type: 'Central' },
  { id: 2, name: 'Central Y', latitude: -22.95, longitude: -43.25, type: 'Central' },
  { id: 2, name: 'Central Y', latitude: -22.95, longitude: -43.25, type: 'Central' },
];

const DELETABLE_TYPES = ['Torre', 'Tower Leader'];

const formCardStyle = {
  display: 'flex',
  flexDirection: 'column',
  gap: '14px',
  background: 'linear-gradient(145deg, #fdfefe 0%, #f3f6ff 100%)',
  padding: '32px',
  borderRadius: '28px',
  width: '360px',
  border: '2px solid #d6e0ff',
  boxShadow: '0 35px 80px rgba(15, 23, 42, 0.18)',
};

function ColorButton({
  color,
  hoverColor,
  textColor = '#fff',
  style = {},
  className = '',
  onMouseEnter,
  onMouseLeave,
  children,
  ...props
}) {
  const [isHovered, setIsHovered] = useState(false);

  const handleMouseEnter = (event) => {
    setIsHovered(true);
    if (onMouseEnter) onMouseEnter(event);
  };

  const handleMouseLeave = (event) => {
    setIsHovered(false);
    if (onMouseLeave) onMouseLeave(event);
  };

  return (
    <button
      {...props}
      className={className}
      onMouseEnter={handleMouseEnter}
      onMouseLeave={handleMouseLeave}
      style={{
        backgroundColor: isHovered && hoverColor ? hoverColor : color,
        color: textColor,
        border: 'none',
        borderRadius: '12px',
        padding: '12px 24px',
        fontWeight: 600,
        cursor: 'pointer',
        boxShadow: '0 12px 25px rgba(15, 23, 42, 0.15)',
        transition: 'background-color 180ms ease, transform 180ms ease',
        ...style,
      }}
    >
      {children}
    </button>
  );
}

function EntityIcon({ src, alt, size = 24 }) {
  return (
    <img
      src={src}
      alt={alt}
      width={size}
      height={size}
      style={{ display: 'block' }}
    />
  );
}

// InfoPopup
function InfoPopup({ item, onConfirm, onCancel }) {
  const canDelete = DELETABLE_TYPES.includes(item?.typeLabel);

  return (
    <div className="absolute bg-white shadow-lg rounded-lg p-4 border border-gray-200 w-60 z-50">
      <h2 className="font-semibold text-lg mb-2 text-gray-900">{item.name}</h2>
      <p className="text-xs text-gray-500 mb-1">ID: {item.id ?? '—'}</p>
      <p className="text-sm text-gray-600">Longitude: {item.longitude ?? '—'}</p>
      <p className="text-sm text-gray-600 mb-3">Latitude: {item.latitude ?? '—'}</p>
      {canDelete && (
        <ColorButton
          onClick={onConfirm}
          color="#ef4444"
          hoverColor="#dc2626"
          style={{ width: '100%', padding: '10px 0' }}
        >
          Excluir
        </ColorButton>
      )}
      <ColorButton
        onClick={onCancel}
        color="#f1f5f9"
        hoverColor="#e2e8f0"
        textColor="#475569"
        style={{
          width: '100%',
          marginTop: canDelete ? '8px' : 0,
          padding: '8px 0',
          boxShadow: 'none',
        }}
      >
        Fechar
      </ColorButton>
    </div>
  );
}

// AddTowerPage (sem alterações)
function AddTowerPage({ onAdd }) {
  const navigate = useNavigate();
  const [form, setForm] = useState({ name: '', latitude: '', longitude: '' });
  const handleChange = (e) => {
    const { name, value } = e.target;
    setForm({ ...form, [name]: value });
  };
  const handleSubmit = (e) => {
    e.preventDefault();
    onAdd({ ...form, id: Date.now(), is_leader: false });
    navigate('/');
  };
  return (
    <div className="p-6 flex flex-col items-center">
      <h1 className="text-2xl font-bold mb-4">Adicionar Torre</h1>
      <ColorButton
        onClick={() => navigate('/')}
        color="#e2e8f0"
        hoverColor="#cbd5f5"
        textColor="#0f172a"
        style={{ alignSelf: 'flex-start', marginBottom: '16px', padding: '8px 18px' }}
      >
        ← Voltar
      </ColorButton>
      <form
        onSubmit={handleSubmit}
        style={formCardStyle}
      >
        <input name="name" placeholder="Nome da Torre" value={form.name} onChange={handleChange} className="form-input" required />
        <input name="latitude" placeholder="Latitude" value={form.latitude} onChange={handleChange} className="form-input" required />
        <input name="longitude" placeholder="Longitude" value={form.longitude} onChange={handleChange} className="form-input" required />
        <ColorButton
          type="submit"
          color="#2563eb"
          hoverColor="#1d4ed8"
          style={{ width: '100%', padding: '10px 0' }}
        >
          Salvar
        </ColorButton>
      </form>
    </div>
  );
}

// AddStructurePage (sem alterações)
function AddStructurePage({ onAdd }) {
  const navigate = useNavigate();
  const [form, setForm] = useState({ name: '', latitude: '', longitude: '', type: 'Platform' });
  const handleChange = (e) => {
    const { name, value } = e.target;
    setForm({ ...form, [name]: value });
  };
  const handleSubmit = (e) => {
    e.preventDefault();
    onAdd({ ...form, id: Date.now() });
    navigate('/');
  };
  return (
    <div className="p-6 flex flex-col items-center">
      <h1 className="text-2xl font-bold mb-4">Adicionar Estrutura</h1>
      <ColorButton
        onClick={() => navigate('/')}
        color="#e2e8f0"
        hoverColor="#cbd5f5"
        textColor="#0f172a"
        style={{ alignSelf: 'flex-start', marginBottom: '16px', padding: '8px 18px' }}
      >
        ← Voltar
      </ColorButton>
      <form onSubmit={handleSubmit} style={formCardStyle}>
        <input name="name" placeholder="Nome da Estrutura" value={form.name} onChange={handleChange} className="form-input" required />
        <input name="latitude" placeholder="Latitude" value={form.latitude} onChange={handleChange} className="form-input" required />
        <input name="longitude" placeholder="Longitude" value={form.longitude} onChange={handleChange} className="form-input" required />
        <select name="type" value={form.type} onChange={handleChange} className="form-input">
          <option value="Platform">Plataforma</option>
          <option value="Central">Central</option>
        </select>
        <ColorButton
          type="submit"
          color="#7c3aed"
          hoverColor="#6d28d9"
          style={{ width: '100%', padding: '10px 0' }}
        >
          Salvar
        </ColorButton>
      </form>
    </div>
  );
}

// AddVehiclePage (sem alterações)
function AddVehiclePage({ onAdd }) {
  const navigate = useNavigate();
  const [form, setForm] = useState({ name: '', type: 'Helicopter', latitude: '', longitude: '' });
  const handleChange = (e) => {
    const { name, value } = e.target;
    setForm({ ...form, [name]: value });
  };
  const handleSubmit = (e) => {
    e.preventDefault();
    onAdd({ ...form, id: Date.now() });
    navigate('/');
  };
  return (
    <div className="p-6 flex flex-col items-center">
      <h1 className="text-2xl font-bold mb-4">Adicionar Veículo</h1>
      <ColorButton
        onClick={() => navigate('/')}
        color="#e2e8f0"
        hoverColor="#cbd5f5"
        textColor="#0f172a"
        style={{ alignSelf: 'flex-start', marginBottom: '16px', padding: '8px 18px' }}
      >
        ← Voltar
      </ColorButton>
      <form onSubmit={handleSubmit} style={formCardStyle}>
        <input name="name" placeholder="Nome do Veículo" value={form.name} onChange={handleChange} className="form-input" required />
        <input name="latitude" placeholder="Latitude" value={form.latitude} onChange={handleChange} className="form-input" required />
        <input name="longitude" placeholder="Longitude" value={form.longitude} onChange={handleChange} className="form-input" required />
        <select name="type" value={form.type} onChange={handleChange} className="form-input">
          <option value="Helicopter">Helicóptero</option>
          <option value="Ship">Navio</option>
        </select>
        <ColorButton
          type="submit"
          color="#16a34a"
          hoverColor="#15803d"
          style={{ width: '100%', padding: '10px 0' }}
        >
          Salvar
        </ColorButton>
      </form>
    </div>
  );
}


function HomePage({ towers, vehicles, structures, setTowers, setVehicles, setStructures }) {
  const [selectedItem, setSelectedItem] = useState(null);
  const [popupPosition, setPopupPosition] = useState({ x: 0, y: 0 });
  const navigate = useNavigate();

  const homeContainerStyle = {
    minHeight: '100vh',
    display: 'flex',
    flexDirection: 'column',
    justifyContent: 'space-between',
    alignItems: 'center',
    padding: '30px 0 24px',
  };

  const mapLayoutStyles = {
    display: 'flex',
    justifyContent: 'space-between',
    alignItems: 'center',
    width: '100%',
    maxWidth: '1400px',
    margin: '0 auto',
    padding: '0 2rem',
  };

  const addButtonConfigs = [
    { key: 'tower', label: '+ Add Tower', to: '/add-tower', color: '#2563eb', hover: '#1d4ed8' },
    { key: 'vehicle', label: '+ Add Vehicle', to: '/add-vehicle', color: '#16a34a', hover: '#15803d' },
    { key: 'structure', label: '+ Add Structure', to: '/add-structure', color: '#7c3aed', hover: '#6d28d9' },
  ];

  const showInfo = (item, e) => {
    const rect = e.currentTarget.getBoundingClientRect();
    const preferredWidth = 260;
    const offset = 14;
    const xPosition =
      item.typeLabel === 'Central'
        ? rect.left - preferredWidth - offset
        : rect.right + offset;
    const yOffset = item.typeLabel === 'Tower Leader' ? rect.top - 40 : rect.top;
    setPopupPosition({ x: xPosition, y: yOffset });
    setSelectedItem(item);
  };

  const handleDelete = () => {
    if (!selectedItem || !DELETABLE_TYPES.includes(selectedItem.typeLabel)) {
      setSelectedItem(null);
      return;
    }

    setTowers((prev) => prev.filter((t) => t.id !== selectedItem.id));
    setSelectedItem(null);
  };

  const platforms = structures.filter((s) => s.type === 'Platform');
  const centrals = structures.filter((s) => s.type === 'Central');
  const towerLeader = towers.find((t) => t.is_leader);
  const regularTowers = towers.filter((t) => !t.is_leader);
  const helicopterVehicles = vehicles.filter((v) => v.type === 'Helicopter');
  const shipVehicles = vehicles.filter((v) => v.type === 'Ship');
  const vehicleColumns = Math.max(helicopterVehicles.length, shipVehicles.length) || 1;

  return (
    <div className="flex flex-col items-center justify-center relative bg-white" style={homeContainerStyle}>
      <h1 className="text-5xl font-extrabold text-gray-900 mb-8">Operational Map</h1>

      {/* Mapa centralizado */}
      <div className="flex flex-row justify-center items-center gap-28 w-full max-w-7xl" style={mapLayoutStyles}>
        {/* Plataformas à esquerda */}
        <div className="grid gap-6 items-center justify-center" style={{ gridTemplateColumns: 'repeat(3, minmax(0, 1fr))' }}>
          {platforms.map((p) => (
            <motion.div
              key={p.id}
              whileHover={{ scale: 1.1 }}
              onClick={(e) => showInfo({ ...p, typeLabel: 'Plataforma' }, e)}
              className="w-16 h-16 bg-purple-500 text-white text-xl font-bold flex items-center justify-center rounded-full cursor-pointer shadow-lg mx-auto"
            >
              <EntityIcon src={platformIcon} alt="Platform icon" size={28} />
            </motion.div>
          ))}
        </div>

        {/* Centro (veículos + torres + líder) */}
        <div className="flex flex-col items-center gap-12">
          {/* Veículos */}
          <div
            className="w-full"
            style={{
              display: 'grid',
              gridTemplateColumns: `repeat(${vehicleColumns}, 70px)`,
              rowGap: '14px',
              columnGap: '14px',
              justifyContent: 'center',
            }}
          >
            {Array.from({ length: vehicleColumns }).map((_, idx) => {
              const vehicle = helicopterVehicles[idx];
              if (!vehicle) {
                return <div key={`helicopter-placeholder-${idx}`} style={{ width: 64, height: 64 }} />;
              }
              return (
                <motion.div
                  key={vehicle.id}
                  whileHover={{ scale: 1.1 }}
                  onClick={(e) => showInfo({ ...vehicle, typeLabel: 'Helicóptero' }, e)}
                  className="flex items-center justify-center w-16 h-16 rounded-full text-white text-xl font-bold cursor-pointer shadow-lg bg-green-500"
                >
                  <EntityIcon src={helicopterIcon} alt="Helicopter icon" size={26} />
                </motion.div>
              );
            })}
            {Array.from({ length: vehicleColumns }).map((_, idx) => {
              const vehicle = shipVehicles[idx];
              if (!vehicle) {
                return <div key={`ship-placeholder-${idx}`} style={{ width: 64, height: 64 }} />;
              }
              return (
                <motion.div
                  key={vehicle.id}
                  whileHover={{ scale: 1.1 }}
                  onClick={(e) => showInfo({ ...vehicle, typeLabel: 'Navio' }, e)}
                  className="flex items-center justify-center w-16 h-16 rounded-full text-white text-xl font-bold cursor-pointer shadow-lg bg-yellow-500"
                >
                  <EntityIcon src={shipIcon} alt="Ship icon" size={26} />
                </motion.div>
              );
            })}
          </div>

          {/* Torres */}
          <div className="bg-blue-50 border-4 border-blue-300 p-10 rounded-3xl grid grid-cols-4 gap-6 shadow-inner">
            {regularTowers.map((t) => (
              <motion.div
                key={t.id}
                whileHover={{ scale: 1.1 }}
                onClick={(e) => showInfo({ ...t, typeLabel: 'Torre' }, e)}
                className="w-14 h-14 bg-blue-500 text-white flex items-center justify-center font-bold rounded-lg cursor-pointer shadow-md"
              >
                <EntityIcon src={towerIcon} alt="Tower icon" />
              </motion.div>
            ))}
          </div>

          {/* Tower Leader */}
          {towerLeader && (
            <motion.div
              key={towerLeader.id}
              whileHover={{ scale: 1.1 }}
              onClick={(e) => showInfo({ ...towerLeader, typeLabel: 'Tower Leader' }, e)}
              className="w-16 h-16 bg-indigo-600 text-white flex items-center justify-center font-bold rounded-lg cursor-pointer shadow-md"
            >
              <EntityIcon src={towerLeaderIcon} alt="Tower leader icon" size={28} />
            </motion.div>
          )}
        </div>

        {/* Centrais à direita */}
        <div className="grid gap-6 items-center justify-center" style={{ gridTemplateColumns: 'repeat(3, minmax(0, 1fr))' }}>
          {centrals.map((c) => (
            <motion.div
              key={c.id}
              whileHover={{ scale: 1.1 }}
              onClick={(e) => showInfo({ ...c, typeLabel: 'Central' }, e)}
              className="w-16 h-16 bg-orange-500 text-white text-xl font-bold flex items-center justify-center rounded-full cursor-pointer shadow-lg mx-auto"
            >
              <EntityIcon src={centralIcon} alt="Central icon" size={28} />
            </motion.div>
          ))}
        </div>
      </div>

      {/* Popup de informações */}
      {selectedItem && (
        <div className="absolute" style={{ top: popupPosition.y, left: popupPosition.x }}>
          <InfoPopup item={selectedItem} onConfirm={handleDelete} onCancel={() => setSelectedItem(null)} />
        </div>
      )}

      {/* Botões de adicionar */}
      <div className="flex gap-6 mt-12">
        {addButtonConfigs.map(({ key, label, to, color, hover }) => (
          <ColorButton
            key={key}
            onClick={() => navigate(to)}
            color={color}
            hoverColor={hover}
            style={{ minWidth: '170px' }}
          >
            {label}
          </ColorButton>
        ))}
      </div>

      {/* Legenda */}
      <div className="absolute bottom-4 right-4 text-sm bg-gray-800 text-white p-4 rounded-lg shadow-md opacity-90 z-40 flex flex-col gap-3">
        {[
          { label: 'Helicóptero', icon: helicopterIcon },
          { label: 'Navio', icon: shipIcon },
          { label: 'Torre', icon: towerIcon },
          { label: 'Tower Leader', icon: towerLeaderIcon },
          { label: 'Plataforma', icon: platformIcon },
          { label: 'Central', icon: centralIcon },
        ].map(({ label, icon }) => (
          <div key={label} className="flex items-center gap-3">
            <EntityIcon src={icon} alt={`${label} icon`} size={20} />
            <span>{label}</span>
          </div>
        ))}
      </div>
    </div>
  );
}


export default function App() {
  const [towers, setTowers] = useState(MOCK_TOWERS);
  const [vehicles, setVehicles] = useState(MOCK_VEHICLES);
  const [structures, setStructures] = useState(MOCK_STRUCTURES);

  const appWrapperStyle = {
    minHeight: '100vh',
    backgroundColor: '#ffffff',
  };

  const contentStyle = {
    width: '100%',
    minHeight: '100vh',
  };

  return (
    <Router>
      <div style={appWrapperStyle}>
        <div style={contentStyle}>
          <Routes>
            <Route
              path="/"
              element={
                <HomePage
                  towers={towers}
                  vehicles={vehicles}
                  structures={structures}
                  setTowers={setTowers}
                  setVehicles={setVehicles}
                  setStructures={setStructures}
                />
              }
            />
            <Route path="/add-tower" element={<AddTowerPage onAdd={(t) => setTowers((prev) => [...prev, t])} />} />
            <Route path="/add-vehicle" element={<AddVehiclePage onAdd={(v) => setVehicles((prev) => [...prev, v])} />} />
            <Route path="/add-structure" element={<AddStructurePage onAdd={(s) => setStructures((prev) => [...prev, s])} />} />
          </Routes>
        </div>
      </div>
    </Router>
  );
}
