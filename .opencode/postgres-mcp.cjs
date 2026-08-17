const { spawn } = require('child_process');
const fs = require('fs');
const os = require('os');
const path = require('path');

const raiz = path.resolve(__dirname, '..');
const archivoEnv = path.join(raiz, '.env');

function cargarEnv(ruta) {
  const variables = {};
  if (!fs.existsSync(ruta)) return variables;
  const lineas = fs.readFileSync(ruta, 'utf8').split(/\r?\n/);
  for (const linea of lineas) {
    const m = linea.match(/^\s*([A-Za-z0-9_]+)\s*=\s*(.*)\s*$/);
    if (m) variables[m[1]] = m[2].replace(/^["']|["']$/g, '');
  }
  return variables;
}

const env = cargarEnv(archivoEnv);
const host = env.BD_HOST || 'localhost';
const puerto = env.BD_PUERTO || '5432';
const usuario = env.BD_USUARIO || 'postgres';
const nombre = env.BD_NOMBRE || 'postgres';
const password = encodeURIComponent(env.BD_PASSWORD || '');
const url = `postgresql://${usuario}:${password}@${host}:${puerto}/${nombre}`;

const paquete = '@modelcontextprotocol/server-postgres';
const entrada = 'dist/index.js';

function buscarEntradaServer() {
  const baseNpx = path.join(os.homedir(), 'AppData', 'Local', 'npm-cache', '_npx');
  if (!fs.existsSync(baseNpx)) return null;
  const candidatos = [];
  for (const carpeta of fs.readdirSync(baseNpx)) {
    const archivo = path.join(baseNpx, carpeta, 'node_modules', paquete, entrada);
    if (fs.existsSync(archivo)) candidatos.push(archivo);
  }
  if (candidatos.length === 0) return null;
  candidatos.sort((a, b) => fs.statSync(b).mtimeMs - fs.statSync(a).mtimeMs);
  return candidatos[0];
}

const entradaServer = buscarEntradaServer();
if (!entradaServer) {
  console.error('postgres-mcp: no se encontro el paquete ' + paquete +
    ' en el cache de npx. Ejecuta una vez: npx -y ' + paquete);
  process.exit(1);
}

const hijo = spawn(process.execPath, [entradaServer, url], {
  stdio: 'inherit',
  env: process.env,
});

hijo.on('error', (err) => {
  console.error('postgres-mcp:', err.message);
  process.exit(1);
});