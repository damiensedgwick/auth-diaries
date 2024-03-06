# Use the official Ubuntu image as the base image
FROM ubuntu:latest

# Set the working directory inside the container
WORKDIR /app

# Copy pre-built binary into the container
COPY bin/auth-diaries /app/auth-diaries

# Copy all static files into the container
COPY static /app/static

# Copy all template files into the container
COPY templates /app/templates

# Copy .env file into the container
COPY .env /app/.env

# Copy the database into the contianer
COPY auth-diaries.db /app/auth-diaries.db

# Expose port 8080 to run the application
EXPOSE 8080

# Command to run the application
CMD ["./auth-diaries"]