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

function drawBbox(ctx, pts) {
	ctx.save()
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
	console.log(xmin, ymin, xmax, ymax)
	ctx.strokeStyle = '#ff0000'
	let m = 10
	ctx.beginPath()
	ctx.rect(xmin - m, ymin - m, xmax - xmin + 2 * m, ymax - ymin + 2 * m)
	ctx.stroke()
	ctx.strokeStyle = '#000000'
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
			drawBbox(ctx, cmd.points)
			break;
		case 'e':
			let [cx, cy, a, b] = cmd.rect
			console.log(cx, cy)
			drawBbox(ctx, [
				[cx - a, cy - b],
				[cx + a, cy + b]
			])
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
	let zoomf = 1;
	let [x, y, w, h] = data.bb.split(',').map(s => Math.floor(parseFloat(s)));
	let c = getCanvas()
	let margin = zoomf * 2;
	c.width = Math.floor(zoomf * w + 2 * margin);
	c.height = Math.floor(zoomf * h + 2 * margin);
	let ctx = c.getContext('2d')
	ctx.translate(margin, margin + h);
	ctx.scale(zoomf, -zoomf);
	for (obj of data['objects']) {
	}
	for (kind of ['objects', 'edges']) {
		for (obj of data[kind]) {
			for (s of ['_draw_', '_ldraw_']) {
				if (s in obj) {
					run(ctx, obj[s])
				}
			}
		}
	}
			
}

async function fetchAndDraw() {
	var data = await get('graph.json')
	draw(data)
}

function main() {
	var c = getCanvas();
	c.addEventListener('mousedown', function(e) {
		var rect = e.target.getBoundingClientRect();
		var x = e.clientX - rect.left; //x position within the element.
		var y = e.clientY - rect.top;  //y position within the element.
		active = findHoverPgon(x, y);
		//redraw();
	}, true);
	c.addEventListener('mouseup', function(e) {
		startX = startY = undefined;
		active = undefined;
		//redraw();
	}, true);

	var startX, startY;
	c.addEventListener('mousemove', function(e) {
		var rect = e.target.getBoundingClientRect();
		var x = e.clientX - rect.left;
		var y = e.clientY - rect.top;
		if (active) {
			if (startX == undefined) {
				startX = x;
				startY = y;
			}
			dx = x - startX;
			dy = y - startY;
			for (let i = 0; i < active.pts.length; i++) {
				active.pts[i][0] += dx;
				active.pts[i][1] += dy;
			}
			startX = x;
			startY = y;
		}
		//redraw();
		//console.log("Active poly:", active);
	}, true);
	fetchAndDraw()
}

main()
