import os
import sqlite3
from flask import Flask, render_template, request, redirect

app = Flask(__name__)
DATABASE = 'passwd.db'
ADMIN_EMAIL = 'huaweilaomke402@gmail.com'
ADMIN_PASSWORD = 'awa114514'


def create_database():
    conn = sqlite3.connect(DATABASE)
    c = conn.cursor()
    c.execute('''CREATE TABLE IF NOT EXISTS users (
                    id INTEGER PRIMARY KEY AUTOINCREMENT,
                    username TEXT NOT NULL,
                    password TEXT NOT NULL)''')
    c.execute('''CREATE TABLE IF NOT EXISTS bases (
                    id INTEGER PRIMARY KEY AUTOINCREMENT,
                    name TEXT NOT NULL,
                    x INTEGER NOT NULL,
                    y INTEGER NOT NULL,
                    z INTEGER NOT NULL,
                    image_path TEXT,
                    owner TEXT)''')
    conn.commit()
    conn.close()


def check_admin(email, password):
    return email == ADMIN_EMAIL and password == ADMIN_PASSWORD


def check_login(username, password):
    conn = sqlite3.connect(DATABASE)
    c = conn.cursor()
    c.execute("SELECT * FROM users WHERE username=? AND password=?", (username, password))
    user = c.fetchone()
    conn.close()
    return user is not None


def get_user_id(username):
    conn = sqlite3.connect(DATABASE)
    c = conn.cursor()
    c.execute("SELECT id FROM users WHERE username=?", (username,))
    user_id = c.fetchone()[0]
    conn.close()
    return user_id


def save_image(image, base_name):
    base_dir = os.path.join(app.root_path, 'images', base_name)
    os.makedirs(base_dir, exist_ok=True)
    image_path = os.path.join(base_dir, image.filename)
    image.save(image_path)
    return image_path


@app.route('/', methods=['GET', 'POST'])
def login():
    if request.method == 'POST':
        username = request.form['username']
        password = request.form['password']
        if check_admin(username, password):
            return redirect('/admin')
        elif check_login(username, password):
            user_id = get_user_id(username)
            return redirect(f'/user/{user_id}')
        else:
            error_message = 'Invalid username or password'
            return render_template('login.html', error_message=error_message)
    return render_template('login.html')


@app.route('/admin', methods=['GET', 'POST'])
def admin():
    if request.method == 'POST':
        base_name = request.form['base_name']
        x = int(request.form['x'])
        y = int(request.form['y'])
        z = int(request.form['z'])
        image = request.files['image']
        owner = request.form['owner']

        image_path = save_image(image, base_name)

        conn = sqlite3.connect(DATABASE)
        c = conn.cursor()
        c.execute("INSERT INTO bases (name, x, y, z, image_path, owner) VALUES (?, ?, ?, ?, ?, ?)",
                  (base_name, x, y, z, image_path, owner))
        conn.commit()
        conn.close()

    conn = sqlite3.connect(DATABASE)
    c = conn.cursor()
    c.execute("SELECT * FROM bases")
    bases = c.fetchall()
    conn.close()
    return render_template('admin.html', bases=bases)


@app.route('/user/<int:user_id>', methods=['GET'])
def user(user_id):
    conn = sqlite3.connect(DATABASE)
    c = conn.cursor()
    c.execute("SELECT * FROM bases WHERE owner=?", (str(user_id),))
    bases = c.fetchall()
    conn.close()
    return render_template('user.html', bases=bases)


if __name__ == '__main__':
    create_database()
    app.run()