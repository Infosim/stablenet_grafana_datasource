/*
 * Copyright: Infosim GmbH & Co. KG Copyright (c) 2000-2020
 * Company: Infosim GmbH & Co. KG,
 *                  Landsteinerstraße 4,
 *                  97074 Wuerzburg, Germany
 *                  www.infosim.net
 */
package util

func FloatPointer(value float64) *float64 {
	result := value
	return &result
}

func StringPointer(value string) *string {
	result := value
	return &result
}
