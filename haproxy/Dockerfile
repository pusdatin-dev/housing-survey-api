FROM haproxy:2.9-alpine

# Switch to root to install packages
USER root

# Install curl (optional, remove if unused)
RUN apk add --no-cache curl

# Copy HAProxy config
COPY haproxy.cfg /usr/local/etc/haproxy/haproxy.cfg

# Revert to the haproxy user for runtime safety
USER haproxy