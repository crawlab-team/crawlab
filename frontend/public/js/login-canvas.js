function initCanvas() {
  let canvas, ctx, circ, nodes, mouse, SENSITIVITY, SIBLINGS_LIMIT, DENSITY, NODES_QTY, ANCHOR_LENGTH, MOUSE_RADIUS,
    TURBULENCE, MOUSE_MOVING_TURBULENCE, MOUSE_ANGLE_TURBULENCE, MOUSE_MOVING_RADIUS, BASE_BRIGHTNESS, RADIUS_DEGRADE,
    SAMPLE_SIZE

  let handle

  // how close next node must be to activate connection (in px)
  // shorter distance == better connection (line width)
  SENSITIVITY = 200
  // note that siblings limit is not 'accurate' as the node can actually have more connections than this value that's because the node accepts sibling nodes with no regard to their current connections this is acceptable because potential fix would not result in significant visual difference
  // more siblings == bigger node
  SIBLINGS_LIMIT = 10
  // default node margin
  DENSITY = 100
  // total number of nodes used (incremented after creation)
  NODES_QTY = 0
  // avoid nodes spreading
  ANCHOR_LENGTH = 100
  // highlight radius
  MOUSE_RADIUS = 200
  // turbulence of randomness
  TURBULENCE = 3
  // turbulence of mouse moving
  MOUSE_MOVING_TURBULENCE = 50
  // turbulence of mouse moving angle
  MOUSE_ANGLE_TURBULENCE = 0.002
  // moving radius of mouse
  MOUSE_MOVING_RADIUS = 600
  // base brightness
  BASE_BRIGHTNESS = 0.12
  // radius degrade
  RADIUS_DEGRADE = 0.4
  // sample size
  SAMPLE_SIZE = 0.5

  circ = 2 * Math.PI
  nodes = []

  canvas = document.querySelector('#canvas')
  if (!canvas) return;
  resizeWindow()
  ctx = canvas.getContext('2d')
  if (!ctx) {
    alert('Ooops! Your browser does not support canvas :\'(')
  }

  function Mouse(x, y) {
    this.anchorX = x
    this.anchorY = y
    this.x = x
    this.y = y - MOUSE_RADIUS / 2
    this.angle = 0
  }

  Mouse.prototype.computePosition = function () {
    // this.x = this.anchorX + MOUSE_MOVING_RADIUS / 2 * Math.sin(this.angle)
    // this.y = this.anchorY - MOUSE_MOVING_RADIUS / 2 * Math.cos(this.angle)
  }

  Mouse.prototype.move = function () {
    let vx = Math.random() * MOUSE_MOVING_TURBULENCE
    let vy = Math.random() * MOUSE_MOVING_TURBULENCE
    if (this.x + vx + MOUSE_RADIUS / 2 > window.innerWidth || this.x + vx - MOUSE_RADIUS / 2 < 0) {
      vx = -vx
    }
    if (this.y + vy + MOUSE_RADIUS / 2 > window.innerHeight || this.y + vy - MOUSE_RADIUS / 2 < 0) {
      vy = -vy
    }
    this.x += vx
    this.y += vy
    // this.angle += Math.random() * MOUSE_ANGLE_TURBULENCE * 2 * Math.PI
    // this.angle -= Math.floor(this.angle / (2 * Math.PI)) * 2 * Math.PI
    // this.computePosition()
  }

  function Node(x, y) {
    this.anchorX = x
    this.anchorY = y
    this.x = Math.random() * (x - (x - ANCHOR_LENGTH)) + (x - ANCHOR_LENGTH)
    this.y = Math.random() * (y - (y - ANCHOR_LENGTH)) + (y - ANCHOR_LENGTH)
    this.vx = Math.random() * TURBULENCE - 1
    this.vy = Math.random() * TURBULENCE - 1
    this.energy = Math.random() * 100
    this.radius = Math.random()
    this.siblings = []
    this.brightness = 0
  }

  Node.prototype.drawNode = function () {
    let color = 'rgba(64, 156, 255, ' + this.brightness + ')'
    ctx.beginPath()
    ctx.arc(this.x, this.y, 2 * this.radius + 2 * this.siblings.length / SIBLINGS_LIMIT / 1.5, 0, circ)
    ctx.fillStyle = color
    ctx.fill()
  }

  Node.prototype.drawConnections = function () {
    for (let i = 0; i < this.siblings.length; i++) {
      let color = 'rgba(64, 156, 255, ' + this.brightness + ')'
      ctx.beginPath()
      ctx.moveTo(this.x, this.y)
      ctx.lineTo(this.siblings[i].x, this.siblings[i].y)
      ctx.lineWidth = 1 - calcDistance(this, this.siblings[i]) / SENSITIVITY
      ctx.strokeStyle = color
      ctx.stroke()
    }
  }

  Node.prototype.moveNode = function () {
    this.energy -= 2
    if (this.energy < 1) {
      this.energy = Math.random() * 100
      if (this.x - this.anchorX < -ANCHOR_LENGTH) {
        this.vx = Math.random() * TURBULENCE
      } else if (this.x - this.anchorX > ANCHOR_LENGTH) {
        this.vx = Math.random() * -TURBULENCE
      } else {
        this.vx = Math.random() * 2 * TURBULENCE - TURBULENCE
      }
      if (this.y - this.anchorY < -ANCHOR_LENGTH) {
        this.vy = Math.random() * TURBULENCE
      } else if (this.y - this.anchorY > ANCHOR_LENGTH) {
        this.vy = Math.random() * -TURBULENCE
      } else {
        this.vy = Math.random() * 2 * TURBULENCE - TURBULENCE
      }
    }
    this.x += this.vx * this.energy / 100
    this.y += this.vy * this.energy / 100
  }

  function Handle() {
    this.isStopped = false
  }

  Handle.prototype.stop = function () {
    this.isStopped = true
  }

  function initNodes() {
    ctx.clearRect(0, 0, canvas.width, canvas.height)
    nodes = []
    for (let i = DENSITY; i < canvas.width; i += DENSITY) {
      for (let j = DENSITY; j < canvas.height; j += DENSITY) {
        nodes.push(new Node(i, j))
        NODES_QTY++
      }
    }
  }

  function initMouse() {
    mouse = new Mouse(canvas.width / 2, canvas.height / 2)
  }

  function initHandle() {
    handle = new Handle()
  }

  function calcDistance(node1, node2) {
    return Math.sqrt(Math.pow(node1.x - node2.x, 2) + (Math.pow(node1.y - node2.y, 2)))
  }

  function findSiblings() {
    let node1, node2, distance
    for (let i = 0; i < NODES_QTY; i++) {
      node1 = nodes[i]
      node1.siblings = []
      for (let j = 0; j < NODES_QTY; j++) {
        node2 = nodes[j]
        if (node1 !== node2) {
          distance = calcDistance(node1, node2)
          if (distance < SENSITIVITY) {
            if (node1.siblings.length < SIBLINGS_LIMIT) {
              node1.siblings.push(node2)
            } else {
              let node_sibling_distance = 0
              let max_distance = 0
              let s
              for (let k = 0; k < SIBLINGS_LIMIT; k++) {
                node_sibling_distance = calcDistance(node1, node1.siblings[k])
                if (node_sibling_distance > max_distance) {
                  max_distance = node_sibling_distance
                  s = k
                }
              }
              if (distance < max_distance) {
                node1.siblings.splice(s, 1)
                node1.siblings.push(node2)
              }
            }
          }
        }
      }
    }
  }

  function redrawScene() {
    if (handle && handle.isStopped) {
      return
    }
    resizeWindow()
    ctx.clearRect(0, 0, canvas.width, canvas.height)
    findSiblings()
    let i, node, distance
    for (i = 0; i < NODES_QTY; i++) {
      node = nodes[i]
      distance = calcDistance({
        x: mouse.x,
        y: mouse.y
      }, node)
      node.brightness = (1 - Math.log(distance / MOUSE_RADIUS * RADIUS_DEGRADE)) * BASE_BRIGHTNESS
    }
    for (i = 0; i < NODES_QTY; i++) {
      node = nodes[i]
      if (node.brightness) {
        node.drawNode()
        node.drawConnections()
      }
      node.moveNode()
    }
    // mouse.move()
    setTimeout(() => {
      requestAnimationFrame(redrawScene)
    }, 50)
  }

  function initHandlers() {
    document.addEventListener('resize', resizeWindow, {passive: true})
    // canvas.addEventListener('mousemove', mousemoveHandler, false)
  }

  function resizeWindow() {
    canvas.width = window.innerWidth
    canvas.height = window.innerHeight
  }

  function mousemoveHandler(e) {
    mouse.x = e.clientX
    mouse.y = e.clientY
  }

  function init() {
    initHandlers()
    initNodes()
    initMouse()
    initHandle()
    redrawScene()
  }

  function reset() {
    handle.isStopped = true
  }

  init()

  window.resetCanvas = reset
}
