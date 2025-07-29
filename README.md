# Sheep Farming Management System

This repository contains a Go backend and a React frontend for managing sheep information such as births, vaccinations and treatments. The project is primarily in Persian.

## Requirements

- **Go 1.23** or newer for the backend
- **Docker** (optional) to run services easily
- **Python 3** (already used in the frontend Dockerfile) for serving the static files

## Running the Backend

The backend exposes a REST API on port `8080`. You can run it locally with Go:

```bash
# from the repository root
cd cmd/api
go run .
```

Alternatively use Docker:

```bash
docker-compose up -d app db
```

This will start the application and a PostgreSQL database.

## Running the Frontend

The `front-react` directory contains a React application which replaces the
legacy static pages in `front`. All functionality such as login, dashboard,
sheep management and vaccination screens has been migrated to React.

```bash
cd front-react
npm install
npm start
```

This will start the development server (usually on `http://localhost:3000`). The
React app expects the backend API on `http://localhost:8080/api/v1` by default.

The legacy static HTML files are still available under the `front` folder and can be served with Python:

```bash
cd front
python3 -m http.server 5500
```

Then open `http://localhost:5500/index.html` in your browser if you need the older UI.

### Parcel

[Parcel](https://parceljs.org/) is a web application bundler that can combine and optimise JavaScript, CSS and assets. It is not required to run this project but can be useful if you want to bundle the frontend into a single output directory. To use it you would create a `package.json`, install Parcel with `npm install parcel`, and run `parcel build index.html`.

### Deploying Frontend to a Server

To deploy, copy the contents of the `front` directory (or Parcel build output) to the web root of your server. Any simple web server such as Nginx, Apache or even the provided Dockerfile can serve the files. For example you could build and run the Docker image:

```bash
cd front
docker build -t sheep-front .
docker run -p 5500:5500 sheep-front
```

## Jalali Calendar

The pages use the local `lib/jalaali.min.js` to convert Jalali dates to Gregorian before sending them to the backend. Make sure the browser can load this script; no external CDN is used.

## Notes on "Cannot use import statement outside a module"

If you see this error in the browser it usually means an ES module was loaded without setting `type="module"` on the `<script>` tag or your bundler is mis‑configured. This project avoids ES module syntax and works without a bundler, so serving the HTML files as shown above should not produce this error.
