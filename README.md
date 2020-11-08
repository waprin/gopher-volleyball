# Gopher Volleyball

Slime volleyball implemented in Go with the SDL 2 bindings

![game demo image](https://github.com/waprin/gopher-volleyball/blob/master/image/slime_demo.png?raw=true)

## About

This is a small project written "just for fun", inspired by Francesc Campoy's  [Flappy Gopher](https://www.youtube.com/watch?v=aYkxFbd6luY)
and reddit user /u/marler8997 []HTMl5 Slime Volleyball](https://www.reddit.com/r/gaming/comments/3b7j47/html5_version_of_slime_volleyball/).

I learned some interesting thing about 2D circle physics while coding this up, check out [my blog post about it](https://waprin.io/2020/10/11/gopher-ball.html).

I didn't make the same effort as some other versions to have total fidelity to the original physics, and I implemented
a basic AI, however both the physics and the AI could use some work. Since this was just a one-off project, I won't 
improve them unless for some reason I'm asked to. 


## Installation

First install the SDL2 bindings for Go by following the instructions
[here](https://github.com/veandco/go-sdl2).

After that, clone this repo and run:

    make run

## Building a binary

To build a binary:

```
    make build 
```    

## Tests

To run tests

```
    make test

```



# License

This project is licensed under the MIT License.