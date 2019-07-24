//   Copyright 2019 MSolution.IO
//
//   Licensed under the Apache License, Version 2.0 (the "License");
//   you may not use this file except in compliance with the License.
//   You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
//   Unless required by applicable law or agreed to in writing, software
//   distributed under the License is distributed on an "AS IS" BASIS,
//   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//   See the License for the specific language governing permissions and
//   limitations under the License.

package ec2

import (
	"regexp"
	"strconv"
	"strings"
)

type InstanceSize struct {
	factor float64
	size   string
	types  []string
}

var (
	INSTANCE = []InstanceSize{
		{1, "nano", []string{"t2", "t3", "t3a"}},
		{2, "micro", []string{"t2", "t1", "t3", "t3a"}},
		{4, "small", []string{"t2", "m1", "t3", "t3a"}},
		{8, "medium", []string{"t2", "m1", "m3", "c1", "t3", "t3a", "a1"}},
		{16, "large", []string{"t2", "m5", "m4", "c5", "c4", "r4", "i3", "m1", "m3", "c3", "r3", "t3", "t3a", "m5d", "m5a", "m5ad", "a1"}},
		{32, "xlarge", []string{"t2", "m5", "m4", "c5", "c4", "p2", "x1e", "r4", "i3", "d2", "m1", "m3", "c1", "c3", "m2", "r3", "i2", "t3", "t3a", "m5d", "m5a", "m5ad", "a1"}},
		{64, "2xlarge", []string{"t2", "m5", "m4", "c5", "c4", "p3", "x1e", "r4", "i3", "h1", "d2", "m3", "c3", "g2", "m2", "r3", "i2", "t3", "t3a", "m5d", "m5a", "m5ad", "a1"}},
		{128, "4xlarge", []string{"m5", "m4", "c5", "c4", "g3", "x1e", "r4", "i3", "h1", "d2", "c3", "m2", "r3", "i2", "m5d", "m5a", "m5ad", "a1"}},
		{256, "8xlarge", []string{"c4", "p2", "p3", "g3", "x1e", "r4", "i3", "h1", "d2", "cc2", "c3", "g2", "cr1", "r3", "i2", "hs1", "m5", "m5d", "m5a"}},
		{288, "9xlarge", []string{"c5"}},
		{320, "10xlarge", []string{"m4"}},
		{384, "12xlarge", []string{"m5", "m5d", "m5a", "m5ad"}},
		{512, "16xlarge", []string{"m4", "p2", "p3", "g3", "x1", "x1e", "r4", "i3", "h1", "m5", "m5d", "m5a"}},
		{576, "18xlarge", []string{"c5"}},
		{768, "24xlarge", []string{"m5", "m5d", "m5a", "m5ad"}},
		{1024, "32xlarge", []string{"x1", "x1e"}},
	}

	rgx = regexp.MustCompile(`([a-zA-Z]+)([\\d])+`)
)

func getEC2RecommendationTypeReason(instance Instance) Recommendation {
	size, family := getInstanceSizeFamily(instance.Type)
	cpuDelta := instance.Stats.Cpu.Average / 0.80
	targetNormFactor := cpuDelta * getNormFactorFromSize(size)
	if instance.Stats.Cpu.Average <= 0 || targetNormFactor == 0 {
		return Recommendation{"", "", getNewGeneration(size, family)}
	}
	recommendedInstance := ""
	finalSize := ""
	var recommendedTemp string
	metaFamily := getSizesForType(family)
	for _, instanceSize := range metaFamily {
		if targetNormFactor <= instanceSize.factor {
			recommendedInstance = family + "." + instanceSize.size
			finalSize = instanceSize.size
			break
		}
		recommendedTemp = instanceSize.size
	}
	if recommendedInstance == instance.Type {
		return Recommendation{"", "", getNewGeneration(size, family)}
	} else if recommendedInstance == "" {
		if recommendedTemp == "" {
			return Recommendation{"", "", getNewGeneration(size, family)}
		}
		return Recommendation{
			InstanceType:  family + "." + recommendedTemp,
			Reason:        getEC2RecommendationReason(getNormFactorFromSize(size), getNormFactorFromSize(recommendedTemp)),
			NewGeneration: getNewGeneration(size, family)}
	}
	reason := getEC2RecommendationReason(getNormFactorFromSize(size), getNormFactorFromSize(finalSize))
	return Recommendation{
		InstanceType:  recommendedInstance,
		Reason:        reason,
		NewGeneration: getNewGeneration(size, family)}
}

func containEc2Type(idx int, family string) bool {
	for _, familyMeta := range INSTANCE[idx].types {
		if familyMeta == family {
			return true
		}
	}
	return false
}

func getInstanceSizeFamily(instanceType string) (size, family string) {
	sizeFamily := strings.Split(instanceType, ".")
	if len(sizeFamily) <= 0 {
		return "", ""
	}
	family = sizeFamily[0]
	if len(sizeFamily) > 1 {
		size = sizeFamily[1]
	}
	return size, family
}

func getSizesForType(currentType string) []InstanceSize {
	size := make([]InstanceSize, 0)
	for idx, value := range INSTANCE {
		if containEc2Type(idx, currentType) {
			size = append(size, value)
		}
	}
	return size
}

func getNormFactorFromSize(size string) float64 {
	for _, instanceSize := range INSTANCE {
		if size == instanceSize.size {
			return instanceSize.factor
		}
	}
	return 0
}

func getEC2RecommendationReason(oldSize, newSize float64) string {
	if oldSize < newSize {
		return "High CPU usage"
	} else if oldSize > newSize {
		return "Low CPU usage"
	}
	return ""
}

func getNewGeneration(size, family string) string {
	for _, instanceSize := range INSTANCE {
		if instanceSize.size == size {
			if newgeneration, available := checkNewGenerationAvailable(size, family, instanceSize); available {
				return strings.Join(newgeneration, ",")
			}
			return ""
		}
	}
	return ""
}

func checkNewGenerationAvailable(size, family string, instanceSize InstanceSize) (recommendedType []string, available bool) {
	available = false
	recommendedType = make([]string, 0)
	actualType := rgx.FindStringSubmatch(family)
	if len(actualType) < 3 {
		return []string{}, false
	}
	actualGen, _ := strconv.Atoi(actualType[2])
	for _, instanceType := range instanceSize.types {
		newGenType := rgx.FindStringSubmatch(instanceType)
		newGen, _ := strconv.Atoi(newGenType[2])
		if len(newGenType) >= 3 && newGenType[1] == actualType[1] && actualGen <= newGen && actualType[0] != newGenType[0] {
			recommendedType = append(recommendedType, instanceType+"."+size)
			available = true
		}
	}
	return recommendedType, available
}
