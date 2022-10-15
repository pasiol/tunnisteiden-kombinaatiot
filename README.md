# Tunnisteiden kombinaatiot

Suodattaa kombinaatiosääntöjen avulla "laittomat" rivit Wilman tuntikirjauksista.

## Postgresql tietokanta

Dockerfile

    FROM postgres:11
    RUN localedef -i fi_FI -c -f UTF-8 -A /usr/share/locale/locale.alias fi_FI.UTF-8
    ENV LANG fi_FI.utf8
    ENV POSTGRES_USER nnnnnnnnnnnnnn
    ENV POSTGRES_PASSWORD nnnnnnnnnnnnnn
    ENV POSTGRES_DB nnnnnnnnnnnnnn

    #!/bin/bash
    docker build -t kombinaatiot .
    docker run -d --rm  -p 5432:5432 --name kombinaatiot_kontti -v ~/data:/var/lib/postgresql/data kombinaatiot

## Alustus

Lue salaisuudet ympäristömuuttujina.

    source set_env.sh

## Käyttö

    ./bin/combinations  -h
    CLI utility to get non valid combinations of the Wilma bookings.

    Usage:
    combinations [command]

    Available Commands:
    completion        Generate the autocompletion script for the specified shell
    getBookings       Importing bookings to pg database.
    getReport         Generates a csv report.
    help              Help about any command
    initializeDb      Initializing db
    updateIdentifiers Updating identifiers

    Flags:
    -h, --help     help for combinations
    -t, --toggle   Help message for toggle

    Use "combinations [command] --help" for more information about a command.

### Tiedokannan alustus

    ./bin/combinations  initializeDb

### Kombinaatiosääntöjen lukeminen

    ./bin/combinations updateIdentifiers --filename ~/lailliset_kombinaatiot.csv

### Opettajien tietojen lukeminen Primuksesta ja Wilma kirjausten lukeminen MongoDB tietokannasta

    ./bin/combinations getBookings

### Raportin tulostaminen

    ./bin/combinations getReport --month 09.2022
 

