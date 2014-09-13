svsetup(1)
==========

Create the skeleton of service of daemontools

Usage
-----

```
$ svsetup <Your Service Name>
```

Then this command generates directories and files, like so;

```
"Your Service Name"/
├── log/
│   ├── main/
│   └── run
└── run
```

And edit `run` files as you like, after create the symbolic link on `/service` to generated directory.

License
-------

MIT

