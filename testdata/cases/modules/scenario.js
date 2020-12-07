function helloWorld(config) {
  return function(context) {
    context.get("/hello")
  }
}

export default helloWorld
