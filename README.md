# Setting up local development:

- ## Server

  - Visit [Download Docker Desktop](https://www.docker.com/products/docker-desktop/) to install docker desktop
  - Run `make start-mongodb` to pull latest mongo docker image and start mongodb on port 27017
  - Visit [Download MongoDB Compass](https://www.mongodb.com/products/compass) to install MongoDB Compass to view and manage data.
  - Open MongoDB compass and click the `Advanced Connection Options`, click the `Authentication` tab and select `Username/Password`. Type `root` in the username box and `pass` in the password tab and hit `Save and Connect`
  - Create a `.env` file in your project's root directory and add `DATABASE_URL=mongodb://root:pass@localhost:27017/?authMechanism=DEFAULT`
  - Run `make start-server` to run the backend server.
  - Run `make test-server` to test the backend server.

# create-svelte

Everything you need to build a Svelte project, powered by [`create-svelte`](https://github.com/sveltejs/kit/tree/master/packages/create-svelte).

## Creating a project

If you're seeing this, you've probably already done this step. Congrats!

```bash
# create a new project in the current directory
npm create svelte@latest

# create a new project in my-app
npm create svelte@latest my-app
```

## Developing

Once you've created a project and installed dependencies with `npm install` (or `pnpm install` or `yarn`), start a development server:

```bash
npm run dev

# or start the server and open the app in a new browser tab
npm run dev -- --open
```

## Building

To create a production version of your app:

```bash
npm run build
```

You can preview the production build with `npm run preview`.

> To deploy your app, you may need to install an [adapter](https://kit.svelte.dev/docs/adapters) for your target environment.
