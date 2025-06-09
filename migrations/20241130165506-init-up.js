module.exports = {
  async up(db, client) {
    await db.createCollection("users");
    await db.createCollection("translations");

    await db.collection("users").createIndex({ "username": 1 }, { unique: true });
  },

  async down(db, client) {
    await db.collection("users").drop();
    await db.collection("translations").drop();
  }
};
