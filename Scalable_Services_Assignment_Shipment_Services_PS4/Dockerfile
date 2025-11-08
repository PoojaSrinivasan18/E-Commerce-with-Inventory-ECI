FROM python:3.11-slim

# Set working directory inside the container
WORKDIR /app

# Copy only requirements first (for efficient Docker caching)
COPY requirements.txt .

# Install dependencies
RUN pip install --no-cache-dir -r requirements.txt

# Copy the rest of the application code
COPY . .

# Expose the port your app runs on
EXPOSE 8001

# Copy and configure entrypoint (wait for DB + db_setup) and make it the container entrypoint
COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

# Default entrypoint will run db setup and start uvicorn
ENTRYPOINT ["/entrypoint.sh"]

