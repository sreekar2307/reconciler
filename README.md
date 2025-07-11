# 📦 Reconciler (MongoDB)

A simple reconciliation system built using **Go** and **MongoDB**, designed to detect mismatched transactions across two collections: `incoming_txns` and `outgoing_txns`.

This project was created as a hands-on exploration of:
- MongoDB's architecture (replica sets, transactions, write concerns, indexing)
- Real-world data modeling and reconciliation logic
- Go MongoDB driver usage

---

## 🚀 Features

- 🔁 Multi-document transactions
- 🔍 Three-way diffing:
  - Present only in incoming
  - Present only in outgoing
  - Present in both but mismatched
- ✅ Mark matched transactions as reconciled
- 💾 MongoDB replica set enabled to support transactions
- 🧪 Sample seed data and index creation

---

## 🧰 Tech Stack

- **Golang**
- **MongoDB (replica set)** via Docker
- **MongoDB Go Driver** (`go.mongodb.org/mongo-driver`)

---

## 🛠️ Setup Instructions

### 1. Clone the Repo

```bash
git clone <your-repo-url>
cd reconciler
```
### 2. Start MongoDB Replica Set

Ensure you have Docker installed, then run:

```bash
docker-compose up -d
docker exec -it mongo mongosh
```
### 3. Initialize the replica set

```js
rs.initiate()
cfg = rs.conf()
cfg.members[0].host = "localhost:27017"
rs.reconfig(cfg, { force: true })
```

### 4. Migrate and Seed Data

```bash
go run github.com/sreekar2307/reconciler/cmd/main migrate
go run github.com/sreekar2307/reconciler/cmd/main seed
```

### 5. Run the Reconciler

```bash
go run github.com/sreekar2307/reconciler/cmd/main recon
```
