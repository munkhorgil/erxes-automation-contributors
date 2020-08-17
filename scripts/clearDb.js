(function() {
  connection = new Mongo();

  const db1 = connection.getDB("erxes");
  const db2 = connection.getDB("erxes_integrations");

  db1.dropDatabase();
  db2.dropDatabase();
})();
