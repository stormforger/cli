definition.session("base", function (session) {
  session.get("/users/configuration", {
    gzip: true,
    tag: "user_configuration",
    headers: {
      "Accept": "application/json",
      "X-DemoApp-Token": session.matchedValue("authenticationToken"),
    },
  });

  session.wait(2, {
    random: true
  });
});
