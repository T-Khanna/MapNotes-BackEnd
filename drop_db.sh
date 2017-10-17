#!/bin/bash

echo "Dropping table"
psql \
  -f drop_db.sql \
  --echo-all \
  postgres

echo "Importing table from Heroku"
heroku pg:pull DATABASE_URL mapnotes_local_db --app mapnotes-backend

