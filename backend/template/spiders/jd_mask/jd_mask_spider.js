const crawlab = require('crawlab-sdk');
const PCR = require('puppeteer-chromium-resolver');

const crawlDetail = async (page, url) => {
	await page.goto(url);
	await page.waitForSelector('#choose-btns');
	await page.waitFor(500);
	
	const hasStock = await page.evaluate(() => {
		return !document.querySelector('.J-notify-stock');
	});
	return hasStock;
};

const crawlPage = async (page) => {
	const items = await page.evaluate(() => {
		const items = [];
		document.querySelectorAll('.gl-item').forEach(el => {
			items.push({ 
				title: el.querySelector('.p-name > a').getAttribute('title'),
				url: 'https:' + el.querySelector('.p-name > a').getAttribute('href'),
			});
		});
		return items;
	});

	for (let i = 0; i < items.length; i++) {
		const item = items[i];
		item['has_stock'] = await crawlDetail(page, item.url);
		await crawlab.saveItem(item);
	}

	await page.waitFor(1000);
};

const main = async () => {
	const pcr = await PCR({
		folderName: '.chromium-browser-snapshots',
		hosts: ["https://storage.googleapis.com", "https://npm.taobao.org/mirrors"],
		retry: 3
	});

	const browser = await pcr.puppeteer.launch({
		headless: true,
		args: ['--no-sandbox'],
		executablePath: pcr.executablePath
	}).catch(function (error) {
		console.log(error);
	});

	const page = await browser.newPage();

	await page.goto('https://www.jd.com/chanpin/270170.html');
	await page.waitForSelector('#J_goodsList');
	await page.waitFor(1000);

	await crawlPage(page);

	while (true) {
		const hasNext = await page.evaluate(() => {
			if (!document.querySelector('.pn-next')) return false
			return !document.querySelector('.pn-next.disabled')
		});

		if (!hasNext) break;

		await page.click('.pn-next');
		await page.waitFor(1000);
		await crawlPage(page);
	}

	await browser.close();
};

(async () => {
	try {
		await main()
	} catch (e) {
		console.error(e)
	}
    
    await crawlab.close();
    // process.exit();
})();