### Example of request
```bash
curl -X POST -F "svg=@<path_to_svg_file.svg>" <hostname>:<port:8080> > output.eps 
```

### How to set up
```bash
docker build -t svg_to_eps .
docker run -d -p 8080:8080 svg_to_eps
```

### Engine
For converting svg to eps is used `inkscape`, the same engine what cloud convert provides.