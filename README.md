# shoring

This was a dumb experiment to look about speed of ensuring a sane user environment.

I believe that there are a host of good unix utilities that can make jobs delightful to solve
on the command line. The problem with the unix philosophy is that it says make a single tool do
one thing well. With no uniform computing environment for unix between MacOS's BSD, Linux, Busybox
and various next generation rust utilities on the perpherial, that do-one-thing approach leads to
a nightmare in trying to get end users to install dependencies to run the joyful tools you write.

We do have some standards on developers laptops:

 - MacOS
 - HomeBrew package manager
 - ASDF version manager

I thought what if we reach into HomeBrew and ASDF which use flat directories on a wicked fast SSD and
check to see what exists, could we write a tool that we could call infront of our bash scripts like:

`shoring --require foo --require "bar~>3.0" --require "woot>=4.1"`

or maybe

`echo {"require":["foo","bar~>3.0","woot>=4.1"]} | shoring`

JSON can fit into an environment variable as well. `export SHORING_DEPS=...`

and have that command return so fast, we didnt care if we called it on every single run.

Sadly I'm a poor golang developer and got distracted but I wanted to commit this anyway.

```
[shoring]> time ./shoring

real	0m0.032s
user	0m0.009s
sys	    0m0.018s
```
