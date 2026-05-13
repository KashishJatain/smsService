# ── Build stage ────────────────────────────────────────────────────────────────
FROM maven:3.9.6-eclipse-temurin-17 AS builder
WORKDIR /app
COPY pom.xml .
# Download dependencies first (layer-cache friendly)
RUN mvn dependency:go-offline -q
COPY src ./src
RUN mvn package -DskipTests -q

# ── Runtime stage ──────────────────────────────────────────────────────────────
FROM eclipse-temurin:17-jre-alpine
WORKDIR /app
RUN addgroup -S sms && adduser -S sms -G sms
COPY --from=builder /app/target/sms-sender-1.0.0.jar app.jar
USER sms
EXPOSE 8080
ENTRYPOINT ["java", "-jar", "app.jar"]
