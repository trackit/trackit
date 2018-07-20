import React, {Component} from 'react';
import * as d3 from 'd3';
import Map from '../../../assets/aws_regions_map.svg';
import PropTypes from "prop-types";

class MapComponent extends Component {

  constructor(props){
    super(props);
    this.createMap = this.createMap.bind(this);
    this.getNodes = this.getNodes.bind(this);
  }

  componentDidMount() {
    this.createMap();
  }

  componentWillReceiveProps(nextProps) {
    if (this.props.data !== nextProps.data)
      this.createMap();
  }

  selectRegion = (region) => {
    this.props.selectRegion(region);
  };

  getNodes() {
    const node = this.node;

    while (node.firstChild)
      node.removeChild(node.firstChild);

    let tooltip = d3.select("div.tooltip#tooltip_map")[0][0];

    if (!tooltip) {
      d3.select("body")
        .append("div")
        .attr("class", "tooltip")
        .attr("id", "tooltip_map")
        .style("opacity", 0);
      tooltip = d3.select("div.tooltip#tooltip_map")[0][0];
    }

    return {node, tooltip};
  }

  createMap() {
    const {node, tooltip} = this.getNodes();

    d3.xml(Map, "image/svg+xml").get((error, map) => {
      if (error)
        node.append(<div className="alert alert-warning" role="alert">Error while getting map ({error})</div>);
      else {
        let importedNode = document.importNode(map.documentElement, true);
        d3.select(importedNode)
          .attr("preserveAspectRatio", "xMidYMid meet")
          .attr("height", 600)
          .attr("width", node.offsetWidth)
          .select("title")
          .html("");
        node.appendChild(importedNode.cloneNode(true));
        Object.keys(this.props.data).forEach((region) => {
          const style = {
            "fill": (this.props.data[region].total ? "#d9534f" : "#f1f1f1"),
            "fill-opacity": (this.props.data[region].total ? this.props.data[region].opacity : 1),
            "cursor": "pointer",
            "pointer-events": "all",
            "stroke": "#777777"
          };
          d3.selectAll("g#AWS-Regions")
            .select("#" + region)
            .on("mouseover", () => {
              tooltip.innerHTML = region + "(" + this.props.data[region].name + ") : <span class='dollar-sign'>$</span>" + parseFloat(this.props.data[region].total.toFixed(2)).toLocaleString();
              d3.select(tooltip)
                .style({
                  opacity: 1,
                  left: (d3.event.pageX + 20) + "px",
                  top: (d3.event.pageY - 30) + "px"
                });
            })
            .on("mouseout", () => {
              tooltip.innerHTML = null;
              d3.select(tooltip)
                .style({opacity: 0});
            })
            .on("click", this.selectRegion.bind(this, region))
            .style(style);
        });
      }
    });
    window.addEventListener("resize", this.createMap);
  };

  render() {
    return (
      <div id="map" ref={node => this.node = node}/>
    );
  }

}

MapComponent.propTypes = {
  data: PropTypes.object.isRequired,
  selectRegion: PropTypes.func.isRequired
};

export default MapComponent;
