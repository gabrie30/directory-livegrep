FROM livegrep/indexer

# Add static files
ADD static/directory-livegrep /root/directory-livegrep
ADD static/entrypoint.sh /root/entrypoint.sh
ADD static/index.html /livegrep/web/templates/index.html
ADD static/layout.html /livegrep/web/templates/common/layout.html

# Run this script on container start
ENTRYPOINT [ "/root/entrypoint.sh" ]