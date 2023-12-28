# Price Whisper
![Coverage](https://img.shields.io/badge/Coverage-82.6%25-brightgreen)
![GitHub release (by tag)](https://img.shields.io/github/v/tag/onyxmoon/hsfl-master-ai-cloud-engineering.svg?sort=semver&label=Version&color=4ccc93d)
[![Run tests (http proxy service)](https://github.com/Onyxmoon/hsfl-master-ai-cloud-engineering/actions/workflows/run-tests-http-proxy-service.yml/badge.svg)](https://github.com/Onyxmoon/hsfl-master-ai-cloud-engineering/actions/workflows/run-tests-http-proxy-service.yml)
[![Run tests (product-service)](https://github.com/Onyxmoon/hsfl-master-ai-cloud-engineering/actions/workflows/run-tests-product-service.yml/badge.svg)](https://github.com/Onyxmoon/hsfl-master-ai-cloud-engineering/actions/workflows/run-tests-product-service.yml)
[![Run tests (shoppinglist-service)](https://github.com/Onyxmoon/hsfl-master-ai-cloud-engineering/actions/workflows/run-tests-shoppinglist-service.yml/badge.svg)](https://github.com/Onyxmoon/hsfl-master-ai-cloud-engineering/actions/workflows/run-tests-shoppinglist-service.yml)
[![Run tests (user-service)](https://github.com/Onyxmoon/hsfl-master-ai-cloud-engineering/actions/workflows/run-tests-user-service.yml/badge.svg)](https://github.com/Onyxmoon/hsfl-master-ai-cloud-engineering/actions/workflows/run-tests-user-service.yml)
[![Run tests (web service)](https://github.com/Onyxmoon/hsfl-master-ai-cloud-engineering/actions/workflows/run-tests-web-service.yml/badge.svg)](https://github.com/Onyxmoon/hsfl-master-ai-cloud-engineering/actions/workflows/run-tests-web-service.yml)

PriceWhisper is your go-to destination for effortless savings on every shopping trip. Compare prices across supermarkets, find the best deals, and save money without the noise. Unlock the power of smart shopping with PriceWhisper.

PriceWhisper is an web application designed to streamline shopping experiences. Users can input their shopping list containing various products, and PriceWhisper meticulously scans multiple retailers, extracting pricing data given merchants.

The web app then analyze this data to provide you with a comprehensive cost breakdown of your shopping list, instantly identifying the most economical merchant for your trip. 

Key features:

- **Price Comparison:** PriceWhisper allows users to compare prices across multiple merchants.
- **Identifying Economical Merchants:** PriceWhisper instantly identifies and highlights the most economical merchant for each item on the user's shopping list and the complete shopping
- **Efficiency:** The platform streamlines the shopping experience, improving efficiency and reducing the effort required to find savings.
- **Data Aggregation:** The database aggregate pricing data from various merchants for one product.

## Prototypical architecture

![Architecture](README.assets/CE_Architecture_Prototype.svg)

## Setup for development
1. Clone repository from Github `git clone git@github.com:Onyxmoon/hsfl-master-ai-cloud-engineering.git`
2. Run `docker compose up -d` to startt the containers defined in the Compose file
3. Install dependencies for frontend `cd src/web-service/frontend` and run `npm ci`
4. Call `npm run build` to build your static files
5. Run `docker compopse down` to stop container, networks, volumes, and images created by up.

## Authors

Dorien Grönwald<br>
dorien.groenwald@stud.hs-flensburg.de<br>
Hochschule Flensburg

Philipp Borucki<br>
philipp.borucki@stud.hs-flensburg.de<br>
Hochschule Flensburg

Jannick Lindner<br>
jannick.lindner@stud.hs-flensburg.de<br>
Hochschule Flensburg

## \<Add more sections here\>
