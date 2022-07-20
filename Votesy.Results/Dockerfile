# So while I'm a fan of the openjdk images, according to their site:
# This image is officially deprecated and all users are recommended to find and use suitable replacements ASAP.
# So I guess I need to do that.

# Also, I started from a base dockerfile here: https://github.com/dockersamples/example-voting-app/blob/master/worker/Dockerfile.j
# But it did not work for various reasons, and this does. 🎉

FROM maven:3.8-openjdk-18-slim AS build

WORKDIR /code

COPY . .
RUN ["mvn", "dependency:resolve"]
RUN ["mvn", "verify"]

# Adding source, compile and package into a fat jar
COPY ["src/main", "/code/src/main"]
RUN ["mvn", "package"]

FROM openjdk:18-slim-buster

COPY --from=build /code/target/results-0.0.1.jar /

CMD ["java", "-jar", "/results-0.0.1.jar"]