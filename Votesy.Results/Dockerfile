FROM maven:3.8-openjdk-18-slim AS build

WORKDIR /code

COPY . .
RUN ["mvn", "dependency:resolve"]
RUN ["mvn", "verify"]

# Adding source, compile and package into a fat jar
COPY ["src/main", "/code/src/main"]
RUN ["mvn", "package"]

FROM openjdk:18-alpine

COPY --from=build /code/target/results-0.0.1.jar /

CMD ["java", "-jar", "/results-0.0.1.jar"]