FROM cgr.dev/chainguard/go AS builder
COPY . /app
RUN cd /app && go build

FROM cgr.dev/chainguard/glibc-dynamic
COPY --from=builder /app/backend .
COPY --from=builder /app/.env .

EXPOSE 8080

ENTRYPOINT ["/backend"]
