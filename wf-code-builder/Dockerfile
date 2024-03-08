# Use an official Node runtime as the base image
FROM node:14-alpine as build

# Set the working directory in the container
WORKDIR /app

# Copy package.json and package-lock.json to the working directory
COPY package*.json ./

# Install dependencies
RUN npm install

# Argument to specify the GitHub URL
ARG GITHUB_URL
# Argument to specify the build command
ARG BUILD_COMMAND
# Clone the repository
RUN git clone $GITHUB_URL .

# Run the build command
RUN $BUILD_COMMAND

# Expose port for running the application
EXPOSE 3000

# Command to start the application
# CMD ["npm", "start"]
