# Liqo Peering Dashboard Client

[React](https://it.reactjs.org/) has been used to write the client. Because of the simplicity of the project, there are few dependencies, and you can see that there isn't a library that helps to handle the state, such as [Redux](https://redux.js.org/)

## Running the code

The client is a standard react application, so you can run it simple using
```bash
npm start
```

While you update the code, you would like to proxy your request to a server. You can learn how to add a proxy by reading the [official documentation](https://create-react-app.dev/docs/proxying-api-requests-in-development/).

## Docker

Once you modify the code you can build a new docker image using the [Dockerfile](./Dockerfile)

```bash
docker build -t <image_name>:<image-tag> .
```
