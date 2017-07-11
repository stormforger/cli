definition.setTarget("testapp.loadtest.party");

definition.setArrivalPhases([
  { duration: 5 * 60, rate: 50, max_users: 2000 },
  { duration: 15 * 60, rate: 60, max_users: 5000 },
]);

definition.setDataSources({
  user: {
    type: "file",
    source: "authentication/users.csv",
    fields: ["id", "email", "password"],
  },
  product: {
    type: "file",
    source: "products/base.csv",
    fields: ["id"],
  },
  likeProbability: {
    type: "random_number",
    range: [1, 100],
  },
});

definition.session("base", function(session) {
  var user = session.pick("user");

  session.post("/users/session", {
    gzip: true,
    tag: "login",
    headers: {
      "Accept": "application/json"
    },
    payload: JSON.stringify({
      "email": user.email(),
      "password": user.password()
    }),
    extraction: {
      jsonpath: {
        "authenticationToken": "user.auth_token",
      }
    },
  });

  session.get("/users/configuration", {
    gzip: true,
    tag: "user_configuration",
    headers: {
      "Accept": "application/json",
      "X-DemoApp-Token": session.matchedValue("authenticationToken"),
    },
  });

  session.wait(2, { random: true });

  session.times(10, function(context) {
    var productId = context.pick("product").id();

    context.get("/products/details/" + productId, {
      gzip: true,
      tag: "product_details",
      headers: {
        "Accept": "application/json",
        "X-DemoApp-Token": session.matchedValue("authenticationToken"),
      },
    });

    // we are going to like a product 30% of the time
    context.if(context.pick("likeProbability"), "lte", 30, function(context) {
      context.wait(10, { random: true });

      context.post("/products/like", {
        gzip: true,
        tag: "product_details",
        headers: {
          "Accept": "application/json",
          "X-DemoApp-Token": context.matchedValue("authenticationToken"),
        },
        payload: JSON.stringify({
          "product_id": productId
        }),
      });
    });

    context.wait(23, { random: true });
  });
});

definition.session("other", function(session) {
  session.forEver(function(context) {
    context.get("/other");
    context.wait(5);
  });
});
