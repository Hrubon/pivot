function getCanvas() {
	return document.getElementById("myCanvas");
}

async function get(url) {
	console.log(`GET ${url}`)
	resp = await fetch(url)
	if (!resp.ok)
		throw `Failed to GET ${url}: ${resp.status}: ${resp.statusText}`
	return resp.json()
}

async function post(url, data) {
	json = JSON.stringify(data)
	console.log(`POST ${url}, payload=${json}`)
	resp = await fetch(url, {
		headers: {
			"Content-Type": "application/json",
		},
		method: "POST",
		body: json,
	})
	if (!resp.ok)
		throw `Failed to POST ${url}: ${resp.status}: ${resp.statusText}`
	return resp.json()
}

function drawPoly(ctx, pts, fill) {
	ctx.beginPath()
	ctx.moveTo(...pts[0])
	for (pt of pts)
		ctx.lineTo(...pt)
	ctx.closePath()
	ctx.strokeStyle = '#000000'
	if (fill)
		ctx.fill()
	else
		ctx.stroke()
}

function drawEllipse(ctx, x, y, a, b) {
	ctx.beginPath()
	ctx.ellipse(x, y, a, b, 0, 0, 2 * Math.PI, false)
	ctx.closePath()
	ctx.stroke()
}

function drawSpline(ctx, pts) {
	ctx.beginPath()
	ctx.moveTo(...pts[0])
	for (pt of pts)
		ctx.lineTo(...pt)
	ctx.stroke()
}

function drawText(ctx, x, y, w, text) {
	ctx.save()
	ctx.translate(x - w/2, y)
	ctx.scale(1, -1)
	ctx.fillText(text, 0, 0)
	ctx.restore()
}

function run(ctx, prog) {
	ctx.save()
	let params = {}
	for (cmd of prog) {
		switch (cmd.op) {
		case 'p':
		case 'P':
			drawPoly(ctx, cmd.points, (cmd.op == 'P'))
			break;
		case 'e':
			let [cx, cy, a, b] = cmd.rect
			drawEllipse(ctx, ...cmd.rect)
			break;
		case 'b':
			drawSpline(ctx, cmd.points)
			break;
		case 'T':
			drawText(ctx, ...cmd.pt, cmd.width, cmd.text)
			break;
		case 'F':
			ctx.font = `${cmd.size}px ${cmd.face}`
			break;
		default:
			/* ignore unknown commands */
		}
	}
	ctx.restore()
}

function draw(data) {
	console.log("Redrawing graph")
	let zoomf = 1;
	let [x, y, w, h] = data.bb.split(',').map(s => Math.floor(parseFloat(s)));
	let c = getCanvas()
	let margin = zoomf * 2;
	c.width = Math.floor(zoomf * w + 2 * margin);
	c.height = Math.floor(zoomf * h + 2 * margin);
	let ctx = c.getContext('2d')
	ctx.translate(margin, margin + h);
	ctx.scale(zoomf, -zoomf);
	for (kind of ['objects', 'edges'])
		for (obj of data[kind])
			for (s of ['_draw_', '_ldraw_'])
				if (s in obj)
					run(ctx, obj[s])
	return trackable(data)
}

function bbox(pts) {
	let m = 10
	let xmin, xmax, ymin, ymax
	xmin = ymin = Number.POSITIVE_INFINITY
	xmax = ymax = Number.NEGATIVE_INFINITY
	for ([x, y] of pts) {
		if (x > xmax)
			xmax = x
		if (x < xmin)
			xmin = x
		if (y > ymax)
			ymax = y
		if (y < ymin)
			ymin = y
	}
	return [xmin - m, ymin - m, xmax + m, ymax + m]
}

function bboxOf(obj) {
	pts = []
	for (cmd of obj['_draw_']) {
		switch (cmd.op) {
		case 'p':
		case 'P':
			pts.push(...cmd.points)
			break
		case 'e':
			let [cx, cy, a, b] = cmd.rect
			pts.push([cx - a, cy - b], [cx + a, cy + b])
			break
		}
	}
	return bbox(pts)
}

function drawTracked(tt) {
	console.log("Redrawing drag 'n' drop data")
	let c = getCanvas()
	let ctx = c.getContext('2d')
	for (t of tt) {
		[xmin, ymin, xmax, ymax] = t[1]
		ctx.save()
		ctx.strokeStyle = '#ff0000'
		ctx.beginPath()
		ctx.rect(xmin, ymin, xmax - xmin, ymax - ymin)
		ctx.stroke()
		ctx.strokeStyle = '#000000'
		ctx.restore()
	}
}

function trackable(data) {
	tt = []
	for (obj of data['objects'])
		if (obj.name.match(/^(net|router)/))
			tt.push([obj, bboxOf(obj)])
	return tt
}

function findInsideObj(tt, x, y) {
	for (t of tt) {
		[xmin, ymin, xmax, ymax] = t[1]
		if (xmin <= x && x <= xmax && ymin <= y && y <= ymax)
			return t[0]
	}
}

function toGvizCoords(x, y) {
	return [x, getCanvas().height - y]
}

async function main() {
	var c = getCanvas();
	let dragged = null
	let tt;
	c.addEventListener('mousedown', function(e) {
		let sx, sy;
		sx = e.clientX
		sy = e.clientY
		dragged = findInsideObj(tt, ...toGvizCoords(sx, sy))
		console.log(`Drag-drop action started at (${sx}, ${sy})`)
	}, true);
	c.addEventListener('mouseup', function(e) {
		let ex = e.clientX
		let ey = e.clientY
		let [gx, gy] = toGvizCoords(ex, ey)
		let [x, y] = dragged.pos.split(",").map(s => parseFloat(s))
		console.log("Orig pos:", x, y)
		x = gx
		y = gy
		post('move', {
			Name: dragged.name,
			X: x,
			Y: y,
		}).then(function(resp) {
			console.log(`Drag-drop action finished, x=${x}, y=${y}`)
			tt = draw(resp.GvizData)
			drawTracked(tt)
		})
	}, true);

	c.addEventListener('mousemove', function(e) {
	}, true);
	var data = await get('graph.json')
	tt = draw(data.GvizData)
	//resp = await post('move', {
	//	Name: "net1",
	//	X: 100,
	//	Y: 0,
	//})
	//resp = await post('move', {
	//	Name: "net0",
	//	X: 188,
	//	Y: 3,
	//})
	//resp = await post('move', {
	//	Name: "router0",
	//	X: 1030,
	//	Y: 430,
	//})
	//tt = draw(resp.GvizData)
	//drawTracked(tt)
}

main()
