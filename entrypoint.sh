#!/bin/sh
# Espera a que la base de datos esté lista
while ! pg_isready -h db -p 5432 -U postgres -d collector -q; do
  echo "Esperando a la base de datos..."
  sleep 1
done
echo "La base de datos está lista.."


ls -l /app/
echo "raiz:"
ls -l /

# Ejecuta la inicialización
/app/migrate-bin db init

# Ejecuta las migraciones
/app/migrate-bin db migrate

# Inicia la aplicación principal
exec /app/main