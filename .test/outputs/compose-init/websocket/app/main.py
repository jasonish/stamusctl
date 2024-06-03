import os
from http.server import BaseHTTPRequestHandler, HTTPServer

class SimpleHTTPRequestHandler(BaseHTTPRequestHandler):
    def _set_headers(self):
        self.send_response(200)
        self.send_header('Content-type', 'text/plain')
        self.end_headers()

    def do_GET(self):
        if self.path == '/ping':
            self._set_headers()
            resp = os.environ.get('RESPONSE', "pong")+'\n'
            self.wfile.write(bytes(resp, 'utf-8'))
        else:
            self.send_error(404, 'Not Found')

def run(server_class=HTTPServer, handler_class=SimpleHTTPRequestHandler):
    port = int(os.environ.get('WEBSOCKET_PORT', 8080))
    server_address = ('', port)
    httpd = server_class(server_address, handler_class)
    print(f'Starting http server on port {port}...')
    httpd.serve_forever()

if __name__ == "__main__":
    run()
