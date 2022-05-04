# Pull the image and call it base
#FROM golang:1.13.7-alpine3.11 as stage1
FROM golang:1.18-alpine3.15 as stage1
# Copy the code
COPY src /codebase/src
RUN ls /codebase/src/main.go
# Build the binary
RUN cd /codebase && go build -v -o /codebase/bin/server ./src/main.go


FROM alpine:3.15 as stage2
# We will copy the final binary from the previous stage to this stage
COPY --from=stage1 /codebase/bin/server /server
ENV PORT=8080
ENV ENVIRONMENT=NOTSMOKE
# Specify the run command for the binary
CMD ["sh", "-c", "/server"]