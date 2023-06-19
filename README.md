# Setting up local development:

- ## Server

  - Visit [Download Docker Desktop](https://www.docker.com/products/docker-desktop/) to install docker desktop
  - Run `make start-mongodb` to pull latest mongo docker image and start mongodb on port 27017
  - Visit [Download MongoDB Compass](https://www.mongodb.com/products/compass) to install MongoDB Compass to view and manage data.
  - Open MongoDB compass and click the `Advanced Connection Options`, click the `Authentication` tab and select `Username/Password`. Type `root` in the username box and `pass` in the password tab and hit `Save and Connect`
  - Create a `.env` file in your project's root directory and add `DATABASE_URL=mongodb://root:pass@localhost:27017/?authMechanism=DEFAULT`
  - Run `make start-server` to run the backend server.
  - Run `make test-server` to test the backend server.
