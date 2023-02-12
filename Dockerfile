FROM mongo:$MONGODB_VERSION

COPY mongod.conf /etc/mongod.conf