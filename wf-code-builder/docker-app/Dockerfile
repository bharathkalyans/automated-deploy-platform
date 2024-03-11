# Use the official Node.js 18 image as the base for this Docker image
FROM node:18

WORKDIR /app

# Copy the entrypoint script from the local filesystem into the Docker image
COPY entrypoint.sh /usr/local/bin/

RUN chmod +x /usr/local/bin/entrypoint.sh

ENTRYPOINT ["entrypoint.sh"]
