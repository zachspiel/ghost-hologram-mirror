# Ghost Hologram Mirror

This repo contains the Go backend and UI to create a hologram mirror similar to the mirror from [Pottery Barn](https://www.potterybarn.com/products/ghost-hologram-mirror/) with additional features. Instead of displaying one static image, this project allows users to select from several different images and gifs to display in realtime.

This project uses a Raspberry PI 4 to display an image behind the mirror. For more information on the setup, refer to the `raspberryPiConfig.md` file in this repo.

## Technology Stack

- GoLang for backend logic and web socket connections
- Tailwind CSS for UI styling

## Running the Project

The project can be started by running the executable located within this project.

```bash
./ghost-hologram-mirror # or ./ghost-hologram-mirror.exe on Windows
```

Alternatively, the server can also be started by running the command `go run .` in the root of this repo.

## Deployment

First build the project, then deploy to fly.io.

- `go build`

- `fly deploy`

## Available Pages

- localhost:8080/
- localhost:8080/display

## Roadmap

- Add controls for images to fade in and out for set interval
- Setup Tailiwnd to not use playground CDN.
