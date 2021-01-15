import requests

def clean_option(option):
    vals = ('n','f','a','d','s')
    while option not in vals:
        option = input('invalid option try again: ')
    return option

def new_note():
    print("new note")

    url = 'http://localhost:8080/api/new'

    title = input('title:')
    note = input ('note:')

    myobj = {"title":title,"note":note}
    response = requests.post(url, data = myobj)

def find_note():
    print("find note")

    url = 'http://localhost:8080/api/find/title'

    title = input('title: ')
    params = {"title":title}
    response = requests.get(url, params=params )

    print("Title", response.json()['title'])
    print("Body", response.json()["note"])


def all_notes():
    url = 'http://localhost:8080/api/all'
    response = requests.post(url)
    json = response.json()
    print("-------ALL NOTES--------")
    for note in json:
        print(" title:",note['title'])
        print('  body:',note['note'])
        print('')
    print("------------------------")


def delete_note():
    print('delete note')

    url = 'http://localhost:8080/api/delete'

    title = input('title:')
    myobj = {"title":title}
    response = requests.post(url, data = myobj)


## PROGRAM LOOP
program_loop = True
while(program_loop):
    print("")
    print("Welcome what you would like do?")
    print(' n- new note')
    print(' f- find note')
    print(' a- all notes')
    print(' d- delete note')
    print(' s- stop program')
    option = clean_option(input('what option will you use: '))
    print("")
    if option == 'n':
        new_note()
    elif option == 'f':
        find_note()
    elif option == 'a':
        all_notes()
    elif option == 'd':
        delete_note()
    elif option == 's':
        program_loop = False
    else:
        print('options are wrong')
