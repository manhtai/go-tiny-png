GoTinyPNG
=========

A simple web service to compress PNG. Use [gopngquant][0].


## Try it:

```sh
docker build -t tinypng .
docker run -d -p 3000:3000 --name tinypng tinypng:latest
```

Visit http://localhost:3000 to compress your PNG images.


## Compare to the "shiny" one:

[TinyPNG][2]




[0]: https://github.com/manhtai/gopngquant
[2]: https://tinypng.com/
