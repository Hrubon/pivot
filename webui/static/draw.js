function canvas() {
	return document.getElementById("myCanvas");
}

var zoomf = 1;
var margin = zoomf * 2;
var c = canvas();
var dragdrop = false;
var ctx = c.getContext("2d");
var width, height;

var active;
c.addEventListener('mousedown', function(e) {
	var rect = e.target.getBoundingClientRect();
	var x = e.clientX - rect.left; //x position within the element.
	var y = e.clientY - rect.top;  //y position within the element.
	active = findHoverPgon(x, y);
	redraw();
}, true);
c.addEventListener('mouseup', function(e) {
	startX = startY = undefined;
	active = undefined;
	redraw();
}, true);

function findHoverPgon(x, y) {
	for (let i = 0; i < pgons.length; i++)
		if (pgons[i].isInside(x, y))
			return pgons[i];
}

var fsize;

class Polygon {
	constructor(pts, fill=false) {
		this.pts = pts;
		this.fill = fill;
		this.acitve = false;
	}

	draw(ctx) {
		ctx.save();
		ctx.beginPath();
		ctx.moveTo(...this.pts[0]);
		ctx.strokeStyle = this.active ? '#ff0000' : '#000000';
		for (let j = 1; j < this.pts.length; j++) {
			ctx.lineTo(...this.pts[j]);
		}
		ctx.closePath();
		if (this.fill)
			ctx.fill();
		else
			ctx.stroke()
		ctx.restore();
	}

	crosses(ux0, uy0, ux1, uy1, vx0, vy0, vx1, vy1) {
		let a = uy1 - uy0;
		let b = ux0 - ux1;
		let c = ux1 * uy0 - ux0 * uy1;
		let d1 = a * vx0 + b * vy0 + c;
		let d2 = a * vx1 + b * vy1 + c;
		return (d1 > 0 && d2 < 0) || (d1 < 0 && d2 > 0);
	}

	isInside(x, y) {
		let intersect = 0;
		let xmin = Number.POSITIVE_INFINITY;
		for (let i = 0; i < this.pts.length; i++) {
			if (this.pts[i][0] < xmin)
				xmin = this.pts[i][0];
		}
		let rx0 = xmin - 2, ry0 = y - 2, rx1 = x, ry1 = y;
		for (let i = 0, j = 1; i < this.pts.length; i++, j++) {
			if (j == this.pts.length)
				j = 0;
			if (this.crosses(rx0, ry0, rx1, ry1, ...this.pts[i], ...this.pts[j])
				&& this.crosses(...this.pts[i], ...this.pts[j], rx0, ry0, rx1, ry1))
				intersect++; // TODO collinearity!
		}
		return intersect % 2 == 1;
	}
}

var pgons = [];

class Ellipse {
	constructor(x, y, a, b) {
		this.x = x;
		this.y = y;
		this.a = a;
		this.b = b;
	}

	draw(ctx) {
		ctx.beginPath();
		ctx.ellipse(this.x, this.y, this.a, this.b, 0, 0, 2 * Math.PI, false);
		ctx.closePath();
		ctx.stroke();
	}
}

class Text {
	constructor(pt, text) {
		this.pt = pt;
		this.text = text;
	}

	draw(ctx) {
		ctx.save();
		ctx.translate(...this.pt);
		ctx.fillText(this.text, 0, 0);
		ctx.restore();
	}
}

class BSpline {
	constructor(pts) {
		this.pts = pts;
	}

	draw(ctx) {
		ctx.beginPath();
		ctx.moveTo(...this.pts[0]);
		for (let j = 1; j < this.pts.length; j++)
			ctx.lineTo(...this.pts[j]);
		ctx.stroke();
	}
}

var elems = [];

function drawit(objs, prop) {
	for (let o = 0; o < objs.length; o++) {
		let obj = objs[o];
		let draw = obj[prop];
		if (!draw)
			return;
		let pos = obj.pos;
		let [x0, y0] = pos.split(',').map(s => parseFloat(s));
		let pts;
		for (let i = 0; i < draw.length; i++) {
			let op = draw[i].op;
			switch (op) {
				case "p":
				case "P":
					pts = draw[i].points;
					for (let k = 0; k < pts.length; k++)
						pts[k][1] = height - pts[k][1];
					poly = new Polygon(pts, op == "P");
					elems.push(poly);
					pgons.push(poly);
					break;

				case "e":
					r = draw[i].rect;
					r[1] = height - r[1];
					ell = new Ellipse(...r);
					elems.push(ell);
					break;

				case "b":
					pts = draw[i].points;
					for (let k = 0; k < pts.length; k++)
						pts[k][1] = height - pts[k][1];
					bsp = new BSpline(pts);
					elems.push(bsp);
					break;
				case "F":
					fsize = Math.floor(draw[i].size)
					ctx.font = fsize + "px " + draw[i].face;
					ctx.textAlign = "center";
					break;
				case "T":
					draw[i].pt[1] = height - draw[i].pt[1];
					txt = new Text(draw[i].pt, draw[i].text);
					elems.push(txt);
					break;
			}
		}
	}
}

function redraw() {
	ctx.clearRect(0, 0, c.width, c.height);
	ctx.moveTo(0, 0);
	for (let i = 0; i < elems.length; i++) {
		elems[i].active = (active == elems[i]);
		elems[i].draw(ctx);
	}
}

async function get(url) {
	return fetch(url).then((r) => r.json())
}

async function all() {
	var data = await get('graph.json')
	console.log(data)

	let bb = data.bb.split(',').map(s => Math.floor(parseFloat(s)));
	width = bb[2], height = bb[3];
	c.width = Math.floor(zoomf * width + 2 * margin);
	c.height = Math.floor(zoomf * height + 2 * margin);
	//ctx.translate(margin, margin);
	ctx.scale(zoomf, zoomf);

	drawit(data.objects, '_draw_');
	drawit(data.objects, '_ldraw_');
	drawit(data.edges, '_draw_');
	drawit(data.edges, '_hdraw_');
	drawit(data.edges, '_ldraw_');

	redraw();
}

all()

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
	redraw();
	//console.log("Active poly:", active);
}, true);
