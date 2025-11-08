# Use an official OpenJDK runtime as a parent image
FROM eclipse-temurin:17-jdk-jammy

WORKDIR /app

# Copy the built jar file into the container
COPY target/notification-service-*.jar app.jar

EXPOSE 8080

ENTRYPOINT ["java", "-jar", "app.jar"]
