
import "./modules/invalid_mod.js"

definition.session("helloworld", function(ctx) {
  ctx.get("/");
});
