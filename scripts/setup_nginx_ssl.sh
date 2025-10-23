#!/bin/bash

set -e

# Color codes for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to print colored output
print_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

# Check if running as root
if [ "$EUID" -ne 0 ]; then 
    print_error "Please run as root (use sudo)"
    exit 1
fi

# Parse command line arguments
while getopts "e:d:p:h" opt; do
    case $opt in
        e) EMAIL="$OPTARG" ;;
        d) DOMAIN="$OPTARG" ;;
        p) PORT="$OPTARG" ;;
        h)
            echo "Usage: $0 -e EMAIL -d DOMAIN -p PORT"
            echo "  -e EMAIL    Email address for Let's Encrypt notifications"
            echo "  -d DOMAIN   Domain name (e.g., example.com)"
            echo "  -p PORT     Backend server port (e.g., 3000)"
            echo "  -h          Show this help message"
            exit 0
            ;;
        \?)
            print_error "Invalid option: -$OPTARG"
            exit 1
            ;;
    esac
done

# Validate required parameters
if [ -z "$EMAIL" ] || [ -z "$DOMAIN" ] || [ -z "$PORT" ]; then
    print_error "Missing required parameters"
    echo "Usage: $0 -e EMAIL -d DOMAIN -p PORT"
    echo "Example: $0 -e admin@example.com -d example.com -p 3000"
    exit 1
fi

print_info "Starting Nginx + Certbot SSL setup..."
print_info "Email: $EMAIL"
print_info "Domain: $DOMAIN"
print_info "Backend Port: $PORT"

# Update system packages
print_info "Updating system packages..."
apt-get update -qq

# Install Nginx
print_info "Installing Nginx..."
apt-get install -y nginx >/dev/null 2>&1

# Install Certbot and Nginx plugin
print_info "Installing Certbot..."
apt-get install -y certbot python3-certbot-nginx >/dev/null 2>&1

# Create Nginx configuration
print_info "Creating Nginx configuration for $DOMAIN..."
cat > /etc/nginx/sites-available/$DOMAIN <<EOF
server {
    listen 80;
    listen [::]:80;
    server_name $DOMAIN www.$DOMAIN;

    location /.well-known/acme-challenge/ {
        root /var/www/html;
    }

    location / {
        proxy_pass http://localhost:$PORT;
        proxy_http_version 1.1;
        proxy_set_header Upgrade \$http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host \$host;
        proxy_set_header X-Real-IP \$remote_addr;
        proxy_set_header X-Forwarded-For \$proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto \$scheme;
        proxy_cache_bypass \$http_upgrade;
    }
}
EOF

# Enable the site
print_info "Enabling Nginx site configuration..."
ln -sf /etc/nginx/sites-available/$DOMAIN /etc/nginx/sites-enabled/

# Remove default site if it exists
if [ -f /etc/nginx/sites-enabled/default ]; then
    rm /etc/nginx/sites-enabled/default
fi

# Test Nginx configuration
print_info "Testing Nginx configuration..."
nginx -t

# Restart Nginx
print_info "Restarting Nginx..."
systemctl restart nginx
systemctl enable nginx

# Allow Nginx through firewall if UFW is active
if command -v ufw >/dev/null 2>&1; then
    if ufw status | grep -q "Status: active"; then
        print_info "Configuring UFW firewall..."
        ufw allow 'Nginx Full' >/dev/null 2>&1
        ufw delete allow 'Nginx HTTP' >/dev/null 2>&1 || true
    fi
fi

# Obtain SSL certificate
print_info "Obtaining SSL certificate from Let's Encrypt..."
certbot --nginx \
    -d $DOMAIN \
    -d www.$DOMAIN \
    --non-interactive \
    --agree-tos \
    --email $EMAIL \
    --redirect \
    --hsts \
    --staple-ocsp

# Test automatic renewal
print_info "Testing certificate renewal..."
certbot renew --dry-run

# Setup automatic renewal (should already be configured by certbot)
print_info "Verifying automatic renewal configuration..."
systemctl status certbot.timer --no-pager || print_warning "Certbot timer not found, manual renewal may be needed"

print_info "Setup complete!"
echo ""
print_info "SSL certificate installed successfully for $DOMAIN"
print_info "Your site should now be accessible at https://$DOMAIN"
print_info "HTTP traffic will automatically redirect to HTTPS"
echo ""
print_info "Certbot will automatically renew certificates before they expire"
print_info "You can manually test renewal with: certbot renew --dry-run"
echo ""
print_info "Nginx configuration: /etc/nginx/sites-available/$DOMAIN"
print_info "SSL certificate location: /etc/letsencrypt/live/$DOMAIN/"
