
# Rebuilds docs and redeploys to s3 bucket


# Add the following to .git/hooks/pre-commit
# 
# SRC_PATTERN="docs/index.asciidoc"
# if git diff --cached --name-only | grep --quiet "$SRC_PATTERN"
# then
#   echo "found docs/index.asciidoc updates"
#   cd /Users/tleyden/Development/gocode/src/github.com/tleyden/cecil/docs/
#   sh redeploy_s3.sh
# fi



# Switch to profile

# Rebuild docs
asciidoctor index.asciidoc

# Push to s3
aws s3 cp . s3://cecil-assets/asciidoc --recursive --exclude "*" --include "*.html" --include "images/**" --profile personal-yahoo
