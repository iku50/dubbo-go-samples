/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package handler

import (
	"context"

	pb "github.com/apache/dubbo-go-samples/online_boutique_demo/recommendationservice/proto/recommendationservice"

	"github.com/apache/dubbo-go-samples/online_boutique_demo/recommendationservice/proto/productservice"
	"github.com/dubbogo/gost/log/logger"
	"golang.org/x/exp/rand"
)

type RecommendationService struct {
	Product productservice.ProductCatalogService
}

const MAX_RESPONSE_COUNT = 5

func (r *RecommendationService) ListRecommendations(ctx context.Context, req *pb.ListRecommendationsRequest) (res *pb.ListRecommendationsResponse, err error) {
	// # fetch list of products from product catalog stub
	catalog, err := r.Product.ListProducts(ctx, &productservice.Empty{})
	if err != nil {
		return nil, err
	}
	filteredProductsIDs := make([]string, 0, len(catalog.Products))
	for _, p := range catalog.Products {
		if contains(p.Id, req.ProductIds) {
			continue
		}
		filteredProductsIDs = append(filteredProductsIDs, p.Id)
	}
	productIDs := sample(filteredProductsIDs, MAX_RESPONSE_COUNT)
	logger.Infof("[Recv ListRecommendations] product_ids=%v", productIDs)
	res = &pb.ListRecommendationsResponse{}
	res.ProductIds = productIDs
	return
}

func contains(target string, source []string) bool {
	for _, s := range source {
		if target == s {
			return true
		}
	}
	return false
}

func sample(source []string, c int) []string {
	n := len(source)
	if n <= c {
		return source
	}
	indices := make([]int, n)
	for i := 0; i < n; i++ {
		indices[i] = i
	}
	for i := n - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		indices[i], indices[j] = indices[j], indices[i]
	}
	result := make([]string, 0, c)
	for i := 0; i < c; i++ {
		result = append(result, source[indices[i]])
	}
	return result
}
