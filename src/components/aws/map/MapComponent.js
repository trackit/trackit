import React, {Component} from 'react';
import * as d3 from 'd3';
import Map from '../../../assets/aws_regions_map.svg';
import PropTypes from "prop-types";

const generateTooltip = (region, data) => {
  const root = document.createElement("div");
  root.classList.add("tooltip-content");

  const header = document.createElement("div");
  header.classList.add("tooltip-header");

  const headerIcon = document.createElement("i");
  headerIcon.classList.add("fa");
  headerIcon.classList.add("fa-globe");

  const headerRegion = document.createElement("div");
  headerRegion.classList.add("tooltip-header-region");

  const headerRegionAWS = document.createElement("div");
  headerRegionAWS.classList.add("tooltip-header-region-aws");
  headerRegionAWS.innerHTML = region;

  const headerRegionPretty = document.createElement("div");
  headerRegionPretty.classList.add("tooltip-header-region-pretty");
  headerRegionPretty.innerHTML = data.name;

  headerRegion.appendChild(headerRegionAWS);
  headerRegion.appendChild(headerRegionPretty);

  header.appendChild(headerIcon);
  header.appendChild(headerRegion);

  const cost = document.createElement("div");
  cost.classList.add("tooltip-cost");

  const costIcon = document.createElement("span");
  costIcon.classList.add("tooltip-cost-icon");
  costIcon.classList.add("dollar-sign");
  costIcon.innerHTML = "$";

  const costValue = document.createElement("div");
  costValue.classList.add("tooltip-cost-value");
  costValue.innerHTML = (data.total < 0.01 && data.total > 0 ? "<0.01" : parseFloat(data.total.toFixed(2)).toLocaleString());

  cost.appendChild(costIcon);
  cost.appendChild(costValue);

  const help = document.createElement("div");
  help.classList.add("tooltip-help");
  help.innerText = "Click for more details";

  root.appendChild(header);
  root.appendChild(cost);
  root.appendChild(help);

  return root;
};

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

    while (node && node.firstChild)
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
    const setupRegion = (region, style, mask=region) => {
      let title;
      switch (region) {
        case "global":
          title = "Global products";
          break;
        case "taxes":
          title = "";
          break;
        case "":
          title = "Other products";
          break;
        default:
          title = region;
      }
      d3.selectAll("g#AWS-Regions")
        .select("#" + mask)
        .on("mouseover", () => {
          tooltip.innerHTML = null;
          tooltip.appendChild(generateTooltip(title, this.props.data[region]));
          d3.select(tooltip)
            .style({
              opacity: 1,
              left: (d3.event.pageX + 10) + "px",
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
    };
    const {node, tooltip} = this.getNodes();

    if (!node)
      return;

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
            "fill": (this.props.data[region].total ? "#4885ed" : "#F2F7FF"),
            "fill-opacity": (this.props.data[region].total ? this.props.data[region].opacity : 1),
            "cursor": "pointer",
            "pointer-events": "all",
            "stroke": "#777777"
          };
          if (region === "global" || region === "") {
            style["stroke"] = "none";
            setupRegion(region, {"cursor": "pointer", "pointer-events": "all"}, "global_toggle");
          } else if (region === "taxes") {
            style["stroke"] = "none";
            setupRegion(region, {"cursor": "pointer", "pointer-events": "all"}, "taxes_toggle");
          }
          setupRegion(region, style, (region !== "" ? region : "global"));
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
