
export let config = {
  url: "https://testapp.loadtest.party",
};

export function helloworld(config) {
  return function(context) {
    context.get(config.url + "/", { tag: "root" });
    context.waitExp(1);
    context.get(config.url + "/data/test.json", {
      tag: "testjson",
    });
  };
};
