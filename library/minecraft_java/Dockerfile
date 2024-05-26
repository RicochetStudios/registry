# Build the wrapper.
FROM golang:1.22.3 as builder

WORKDIR /go/src/ricochet/minecraft_java
COPY . /go/src/ricochet/minecraft_java
RUN go build -o wrapper .

# Build the final image.
# Uses the common Minecraft server image as the base.
FROM itzg/minecraft-server:stable

# Copy the wrapper from the builder stage.
COPY --from=builder /go/src/ricochet/minecraft_java/wrapper /wrapper
RUN chmod +x wrapper

ENTRYPOINT /wrapper
HEALTHCHECK --start-period=1m --interval=5s --retries=24 CMD /wrapper -health