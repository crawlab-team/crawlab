# Build stage
FROM node:20-alpine AS build

# Set the working directory in the container
WORKDIR /app

# Install pnpm
RUN npm install -g pnpm

# Copy application code
ADD . .

# Install project dependencies
RUN pnpm install

# Build the application
RUN pnpm run build

# Production stage
FROM alpine:3.14

# Copy only the built artifacts from the build stage
COPY --from=build /app/dist /app/dist
