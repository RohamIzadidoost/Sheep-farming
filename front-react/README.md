# React Frontend

This directory contains the React implementation of the full frontend. It replaces the old static pages under `../front`.

## Development

```bash
npm install
npm start
```

The app expects the backend to be available on `http://localhost:8080/api/v1`. You can override this by setting `REACT_APP_API_BASE` when starting the server.

## Production build

To create a static build run:

```bash
npm run build
```

The output will be generated in the `build` directory.
